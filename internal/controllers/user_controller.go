package controllers

import (
    "central-auth/internal/db"
)

type User struct {
    ID string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    LastIP string `json:"lastIp"`
}

type UserController struct {
    db *database.Database
}

func NewUserController(db *database.Database) *UserController {
    return &UserController{
        db: db,
    }
}
