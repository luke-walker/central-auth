package server

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "central-auth/internal/db"
)

type AuthServer struct {
    addr string
    db *database.Database
    router chi.Router
}

func NewAuthServer(addr string, dbURL string) (*AuthServer, error) {
    db, err := database.NewDatabase(dbURL)
    if err != nil {
        return nil, err
    }

    router := chi.NewRouter()
    // Endpoints
    // ...

    return &AuthServer{
        addr: addr,
        db: db,
        router: router,
    }, nil
}

func (server *AuthServer) Start() {
    http.ListenAndServe(server.addr, server.router)
}

func (server *AuthServer) Close() {
    server.db.Close()
}
