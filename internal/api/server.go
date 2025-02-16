package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(serverPort string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%s", serverPort),
			Handler: handler,
		},
	}
}
