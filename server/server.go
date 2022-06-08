package server

import (
	"bytes"
	"errors"
	"io"
)

var errBufTooSmall = errors.New("buffer is too small to fit a sigle message")

// InMemory store all data in memory
type InMemory struct {
	buf []byte
}

// Ack mark the current chunk as done and delete it's contents.
func (s *InMemory) Ack() error {
	s.buf = nil
	return nil
}

// Write accepts the messages from the clients and store them.
func (s *InMemory) Write(msgs []byte) error {
	s.buf = append(s.buf, msgs...)
	return nil
}

// Read copies the data from the in-memory store and writes
// the data read to the the provided Writer, starting with
// the offset provided.
func (s *InMemory) Read(off uint64, maxSize uint64, w io.Writer) error {
	maxOff := uint64(len(s.buf))
	if off >= maxOff {
		return nil
	} else if off+maxSize > maxOff {
		w.Write(s.buf[off:])
		return nil
	}

	truncated, _, err := cutToLastMessage(s.buf[off : off+maxSize])
	if err != nil {
		return err
	}

	if _, err = w.Write(truncated); err != nil {
		return err
	}

	return nil
}

func cutToLastMessage(res []byte) (truncated []byte, rest []byte, err error) {
	n := len(res)

	if n == 0 {
		return res, nil, nil
	}

	if res[n-1] == '\n' {
		return res, nil, nil
	}

	lastPos := bytes.LastIndexByte(res, '\n')
	if lastPos < 0 {
		return nil, nil, errBufTooSmall
	}

	return res[:lastPos+1], res[lastPos+1:], nil
}
