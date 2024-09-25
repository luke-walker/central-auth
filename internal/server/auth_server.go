package server

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/luke-walker/go-validate"

    "central-auth/internal/controllers"
    "central-auth/internal/db"
    "central-auth/pkg/middleware" // beware importing /internal/middleware
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
    authController := controllers.NewAuthController(db)
    serverController := controllers.NewServerController(db)

    /* Routers */
    r := chi.NewRouter()
    r.Route("/auth", func(r chi.Router) {
        /* /auth */
        r.With(middleware.AuthenticateUser(db, false)).Post("/", func(w http.ResponseWriter, r *http.Request) {}) // could this function be nil instead?

        /* /auth/admin */
        r.With(middleware.AuthenticateUser(db, true)).Post("/admin", func(w http.ResponseWriter, r *http.Request) {})

        /* /auth/login/{serverToken} */
        r.Get("/login/{serverToken}", authController.GetLoginPage)
        r.With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "username": { "required": true },
                "password": { "required": true },
            },
        }.ValidateData).Post("/login/{serverToken}", authController.AttemptUserLogin)

        /* /auth/signup/{serverToken} */
        r.Get("/signup/{serverToken}", authController.GetSignupPage)
        r.With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "username": { "required": true },
                "password": { "required": true },
            },
        }.ValidateData).Post("/signup/{serverToken}", authController.AttemptUserSignUp)
    })
    r.Route("/server", func(r chi.Router) {
        /* /server */
        r.With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "name": { "required": true },
                "addrs": { "required": true },
                "redirect": { "required": true },
            },
        }.ValidateData).Post("/", serverController.CreateServer)
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
