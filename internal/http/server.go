package http

import (
	"log"
	"net/http"
)

type Server struct {
	addr string
	// router  + app.Service
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run() error {
	log.Printf("starting http server on %s", s.addr)
	// router по openapi.yml
	return http.ListenAndServe(s.addr, nil)
}
