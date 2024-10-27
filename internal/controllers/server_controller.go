package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/luke-walker/central-auth/internal/db"
)

type Server struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Addrs []string `json:"addrs"`
    Proxy string `json:"proxy"`
    Redirect string `json:"redirect"`
    Token string `json:"token"`
}

type ServerController struct {
    db *database.Database
}

func NewServerController(db *database.Database) *ServerController {
    return &ServerController{
        db: db,
    }
}

func (c *ServerController) CreateServer(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errInvalidRequestMethod(w)
        return
    }

    var data Server
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := c.db.CreateServer(data.Name, data.Addrs, data.Proxy, data.Redirect); err != nil {
        http.Error(w, fmt.Sprintf("Error creating server: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
