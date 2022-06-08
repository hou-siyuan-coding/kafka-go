package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

const defaultScratchSizew = 64 * 1024

type Simple struct {
	addrs []string
	cl    *http.Client
	off   uint64
}

func NewSimple(addrs []string) *Simple {
	return &Simple{
		addrs: addrs,
		cl:    &http.Client{},
	}
}

func (s *Simple) Send(msgs []byte) error {
	resp, err := s.cl.Post(s.addrs[0]+"/write", "application/octet-stream", bytes.NewReader(msgs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		return fmt.Errorf("sending data: http code %d, %s", resp.StatusCode, b.String())
	}

	io.Copy(ioutil.Discard, resp.Body)
	return nil
}

func (s *Simple) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchSizew)
	}

	addrIdx := rand.Intn(len(s.addrs))
	addr := s.addrs[addrIdx]
	readUrl := fmt.Sprintf("%s/read?off=%d&maxSize=%d", addr, s.off, len(scratch))

	resp, err := s.cl.Get(readUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		return nil, fmt.Errorf("get data: http code %d, %s", resp.StatusCode, b.String())
	}

	b := bytes.NewBuffer(scratch[0:0])
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return nil, err
	}

	// 0 bytes read but no errors means the end of file by convention.
	if b.Len() == 0 {
		log.Println("0 bytes read but no errors means the end of file by convention.")
		if err := s.ackCurrentChunk(addr); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	s.off += uint64(b.Len())
	return b.Bytes(), nil
}

func (s *Simple) ackCurrentChunk(addr string) error {
	resp, err := s.cl.Get(addr + "/ack")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		return fmt.Errorf("ack: http code %d, %s", resp.StatusCode, b.String())
	}

	io.Copy(ioutil.Discard, resp.Body)
	return nil
}
