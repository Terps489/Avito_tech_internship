package http

import (
	"log"
	"net/http"

	"github.com/terps489/avito_tech_internship/internal/app"
)

type Server struct {
	addr    string
	service *app.Service
}

func NewServer(addr string, svc *app.Service) *Server {
	return &Server{
		addr:    addr,
		service: svc,
	}
}

func (s *Server) Run() error {
	log.Printf("starting http server on %s", s.addr)

	// TODO: implement HTTP handlers
	mux := http.NewServeMux()

	return http.ListenAndServe(s.addr, mux)
}
