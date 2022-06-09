package server

import (
	"io"
	"os"
)

// OnDisk store all data in Disk
type OnDisk struct {
	fp *os.File
}

func NewOnDisk(fp *os.File) *OnDisk {
	return &OnDisk{fp: fp}
}

// Ack mark the current chunk as done and delete it's contents.
func (s *OnDisk) Ack() error {
	if err := s.fp.Truncate(0); err != nil {
		return nil
	}

	// truncate remove file contents but don't change offset
	// set it to begin
	_, err := s.fp.Seek(0, 0)
	return err
}

// Write accepts the messages from the clients and store them.
func (s *OnDisk) Write(msgs []byte) error {
	_, err := s.fp.Write(msgs)
	return err
}

// Read copies the data from the in-memory store and writes
// the data read to the the provided Writer, starting with
// the offset provided.
func (s *OnDisk) Read(off uint64, maxSize uint64, w io.Writer) error {
	buf := make([]byte, maxSize)
	n, err := s.fp.ReadAt(buf, int64(off))
	if n == 0 {
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}

	// Read until the last message.
	// Don't send the incomplete part of the last
	// message if it is cut in half.
	truncated, _, err := cutToLastMessage(buf[0:n])
	if err != nil {
		return err
	}

	if _, err := w.Write(truncated); err != nil {
		return err
	}

	return nil
}
