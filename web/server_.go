package web

import (
	"io"
	"net/http"

	"github.com/hou-siyuan-coding/kafka-go/server"
)

type Server_ struct {
	s *server.InMemory
}

func (s *Server_) handler(w http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	case "/write":
		io.WriteString(w, "hello write!")
	default:
		io.WriteString(w, "hello world")
	}
}

func (s *Server_) Serve() {
	http.HandleFunc("/", s.handler)
	http.ListenAndServe(":8080", nil)
}
