package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"

    "github.com/joho/godotenv"

    "github.com/luke-walker/central-auth/internal/server"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("error loading .env file: %v", err)
    }

    serverHost := os.Getenv("SERVER_HOST")
    serverPort := os.Getenv("SERVER_PORT")
    serverAddr := fmt.Sprintf("%s:%s", serverHost, serverPort)

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

    authServer, err := server.NewAuthServer(serverAddr, dbURL)
    if err != nil {
        log.Fatalf("error starting authentication server: %v", err)
    }
    fmt.Println("starting authentication server")
    go authServer.Start()
    defer authServer.Close()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    <-stop

    fmt.Println("stopping authentication server")
}
