package main

import (
	"net/http"

	"github.com/mrchip53/go-a-garden/dist"
)

// Use go1.22 net/http package for router

type server struct {
	server *http.Server
	router *http.ServeMux
}

func newServer() *server {
	s := &server{}

	s.router = http.NewServeMux()

	s.router.Handle("/", http.FileServer(http.FS(dist.WebAssets)))
	s.router.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:    ":8001",
		Handler: s.router,
	}

	return s
}

func (s *server) GET(path string, handler http.HandlerFunc) {
	s.router.HandleFunc("GET "+path, handler)
}

func (s *server) POST(path string, handler http.HandlerFunc) {
	s.router.HandleFunc("POST "+path, handler)
}

func (s *server) PUT(path string, handler http.HandlerFunc) {
	s.router.HandleFunc("PUT "+path, handler)
}

func (s *server) DELETE(path string, handler http.HandlerFunc) {
	s.router.HandleFunc("DELETE "+path, handler)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
