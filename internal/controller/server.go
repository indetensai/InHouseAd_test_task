package controller

import (
	"net/http"
	"task/internal/config"
	"task/internal/repository"
	"task/internal/service"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router  chi.Router
	Service service.CheckService
	Body    http.Server
}

func NewServer(repo repository.Repository, cfg config.Config, service service.CheckService) Server {
	var s Server
	s.Router = chi.NewRouter()
	s.Service = service
	s.Body = http.Server{Addr: cfg.ListenAddress, Handler: s.Router}
	AddHandlers(s.Router, s.Service, repo)
	return s
}
