package main

import (
	"io"
	"net/http"

	"github.com/hou-siyuan-coding/kafka-go/server"
	"github.com/hou-siyuan-coding/kafka-go/web"
)

type Server_ struct {
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

func main() {
	// fasthttp.ListenAndServe(":8080", HTTPHandler)
	web.NewServer(&server.InMemory{}).Serve()
	// s := Server_{}
	// s.Serve()
}
