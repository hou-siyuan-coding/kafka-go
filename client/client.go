package client

import (
	"bytes"
	"errors"
)

var errBufTooSmall = errors.New("buffer is too small to fit a sigle message")

const defaultScratchSizew = 64 * 1024

type Simple struct {
	addrs []string

	buf     bytes.Buffer
	restBuf bytes.Buffer
}

func NewSimple(addrs []string) *Simple {
	return &Simple{
		addrs: addrs,
	}
}

func (s *Simple) Send(msgs []byte) error {
	_, err := s.buf.Write(msgs)
	return err
}

func (s *Simple) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchSizew)
	}

	startOff := 0

	if s.restBuf.Len() > 0 {
		if s.restBuf.Len() >= len(scratch) {
			return nil, errBufTooSmall
		}

		n, err := s.restBuf.Read(scratch)
		if err != nil {
			return nil, err
		}

		s.restBuf.Reset()
		startOff += n
	}

	n, err := s.buf.Read(scratch[startOff:])
	if err != nil {
		return nil, err
	}

	truncated, rest, err := cutToLastMessage(scratch[:startOff+n])
	if err != nil {
		return nil, err
	}

	s.restBuf.Write(rest)

	return truncated, nil
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
