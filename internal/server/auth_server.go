package server

import (
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/go-chi/httprate"
    _ "github.com/joho/godotenv/autoload"
    "github.com/luke-walker/go-validate"

    "central-auth/internal/controllers"
    "central-auth/internal/db"
    "central-auth/pkg/middleware"
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
    userController := controllers.NewUserController(db)

    /* Routers */
    r := chi.NewRouter()
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
        AllowCredentials: true,
    }))
    r.Use(httprate.LimitByIP(100, time.Minute))
    r.Use(chiMiddleware.RealIP)

    r.Route("/auth", func(r chi.Router) {
        r.Group(func(r chi.Router) {
            r.Use(middleware.AddBearerTokenHeader)

            /* /auth */
            r.With(middleware.AuthenticateUser(db, false)).Get("/", func(w http.ResponseWriter, r *http.Request) {}) // could this function be nil instead?

            /* /auth/admin */
            r.With(middleware.AuthenticateUser(db, true)).Get("/admin", func(w http.ResponseWriter, r *http.Request) {})

            /* /auth/logout */
            r.Post("/logout", authController.AttemptUserLogout)
        })

        /* /auth/login/{serverToken} */
        r.Get("/login/{serverToken}", authController.GetLoginPage)
        r.With(httprate.LimitByIP(5, time.Minute)).With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "username": { "required": true },
                "password": { "required": true },
            },
        }.ValidateData).Post("/login/{serverToken}", authController.AttemptUserLogin)

        /* /auth/signup/{serverToken} */
        r.Get("/signup/{serverToken}", authController.GetSignupPage)
        r.With(httprate.LimitByIP(3, 5*time.Minute)).With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "username": { "required": true },
                "password": { "required": true },
            },
        }.ValidateData).Post("/signup/{serverToken}", authController.AttemptUserSignUp)

    })
    r.Route("/server", func(r chi.Router) {
        r.Use(middleware.AuthenticateUser(db, true))

        /* /server */
        r.With(govalidate.Validator{
            Fields: govalidate.FieldsMap{
                "name": { "required": true },
                "addrs": { "required": true },
                "redirect": { "required": true },
            },
        }.ValidateData).Post("/", serverController.CreateServer)
    })
    r.Route("/user", func(r chi.Router) {
        r.Use(middleware.AddBearerTokenHeader)

        /* /user */
        r.Get("/", userController.GetUserInfo)
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
