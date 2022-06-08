package web

import (
	"fmt"
	"io"

	"github.com/hou-siyuan-coding/kafka-go/server"
	"github.com/valyala/fasthttp"
)

type Server struct {
	s *server.InMemory
}

func NewServer(s *server.InMemory) *Server {
	return &Server{s: s}
}

func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/write":
		s.writeHandler(ctx)
	case "/read":
		s.readHandler(ctx)
	case "/ack":
		s.ackHandler(ctx)
	default:
		ctx.WriteString("hello world!")
	}
}

func (s *Server) writeHandler(ctx *fasthttp.RequestCtx) {
	if err := s.s.Write(ctx.Request.Body()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
	}
}

func (s *Server) readHandler(ctx *fasthttp.RequestCtx) {
	off, err := ctx.QueryArgs().GetUint("off")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("bad `off` GET param: %v", err))
		return
	}

	maxSize, err := ctx.QueryArgs().GetUint("maxSize")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("bad `maxSize` GET param: %v", err))
		return
	}

	err = s.s.Read(uint64(off), uint64(maxSize), ctx)
	if err != nil && err != io.EOF {
		ctx.SetStatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
}

func (s *Server) ackHandler(ctx *fasthttp.RequestCtx) {
	if err := s.s.Ack(); err != nil {
		ctx.SetStatusCode(500)
		ctx.WriteString(err.Error())
	}
}

func (s *Server) Serve() error {
	return fasthttp.ListenAndServe(":8080", s.handler)
}