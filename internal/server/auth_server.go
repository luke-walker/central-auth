package server

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/luke-walker/go-validate"

    "central-auth/internal/controllers"
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

    /* Controllers */
    serverController := controllers.NewServerController(db)

    /* Routers */
    r := chi.NewRouter()
    r.Route("/servers", func(r chi.Router) {
        r.With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "name": { "required": true },
                "addrs": { "required": true },
                "redirect": { "required": false },
            },
        }.ValidateJSON).Post("/", serverController.CreateServer)
    })

    return &AuthServer{
        addr: addr,
        db: db,
        router: r,
    }, nil
}

func (authServer *AuthServer) Start() {
    http.ListenAndServe(authServer.addr, authServer.router)
}

func (authServer *AuthServer) Close() {
    authServer.db.Close()
}
