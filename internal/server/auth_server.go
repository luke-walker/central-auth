package server

import (
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
)

type AuthServer struct {
    address string
    router chi.Router
}

func NewAuthServer(host string, port int) *AuthServer {
    router := chi.NewRouter()

    // Endpoints
    // ...

    return &AuthServer{
        address: fmt.Sprintf("%s:%d", host, port),
        router: router,
    }
}

func (server *AuthServer) Start() {
    http.ListenAndServe(server.address, server.router)
}
