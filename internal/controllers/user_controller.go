package controllers

import (
    "encoding/json"
    "net/http"

    "central-auth/internal/db"
)

/*type User struct {
    ID string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    LastIP string `json:"lastIp"`
    Admin bool `json:"admin"`
}*/

type UserController struct {
    db *database.Database
}

func NewUserController(db *database.Database) *UserController {
    return &UserController{
        db: db,
    }
}

func (c *UserController) GetUserInfo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        errInvalidRequestMethod(w)
        return
    }

    accessToken, err := GetBearerToken(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userInfo, numRows, err := c.db.GetUserInfoByAccessToken(accessToken)
    if err != nil {
        http.Error(w, "Error retrieving user information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, "Could not find matching session", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(struct{
        Token string `json:"token"`
        Username string `json:"username"`
        LastIP string `json:"last_ip"`
    }{
        Token: userInfo.Token,
        Username: userInfo.Username,
        LastIP: userInfo.LastIP,
    })
}
