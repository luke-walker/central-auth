package controllers

import (
    "fmt"
    "html/template"
    "net/http"
    "os"
    "time"

    "github.com/go-chi/chi/v5"
    _ "github.com/joho/godotenv/autoload"

    "central-auth/internal/crypto"
    "central-auth/internal/db"
)

type AuthController struct {
    db *database.Database
}

func NewAuthController(db *database.Database) *AuthController {
    return &AuthController{
        db: db,
    }
}

func (c *AuthController) GetLoginPage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        errInvalidRequestMethod(w)
        return
    }

    /* Verify Server */
    serverToken := chi.URLParam(r, "serverToken")
    serverInfo, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, fmt.Sprintf("Could not find matching server token '%s'", serverToken), http.StatusUnauthorized)
        return
    }

    /* Send login page */
    loginTemplate, err := template.ParseFiles(fmt.Sprintf("%s/internal/templates/login.html.tmpl", os.Getenv("PROJECT_PATH")))
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing login.html template: %v", err), http.StatusInternalServerError)
        return
    }
    loginTemplate.Execute(w, struct{
        ServerName string
        ServerToken string
    }{
        ServerName: serverInfo.Name,
        ServerToken: chi.URLParam(r, "serverToken"),
    })
}

func (c *AuthController) AttemptUserLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errInvalidRequestMethod(w)
        return
    }
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing URL-encoded form", http.StatusInternalServerError)
        return
    }

    /* User IP */
    userIp := r.Header.Get("X-Forwarded-For")
    if userIp == "" {
        userIp = r.RemoteAddr
    }

    /* Verify Server */
    serverToken := chi.URLParam(r, "serverToken")
    serverInfo, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, fmt.Sprintf("Could not find matching server token '%s'", serverToken), http.StatusUnauthorized)
        return
    }

    /* Verify User */
    username := r.PostForm["username"][0]
    password := r.PostForm["password"][0]
    userInfo, numRows, err := c.db.GetUserInfoByUsername(username)
    if err != nil {
        http.Error(w, "Error retrieving user information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, fmt.Sprintf("Could not find matching username '%s'", username), http.StatusUnauthorized)
        return
    }
    if err = crypto.CheckPasswordHash(password, userInfo.Password); err != nil {
        http.Error(w, fmt.Sprintf("Incorrect password for username '%s'", username), http.StatusUnauthorized)
        return
    }

    /* Check if session already exists -- and delete */
    if _, err = c.db.DestroySession(userIp, userInfo.Token); err != nil {
        http.Error(w, "Error attempting to destroy any pre-existing session", http.StatusInternalServerError)
        return
    }

    /* Create Session */
    const EXPIRES_DAYS = 3
    expiresTime := time.Now().Add(EXPIRES_DAYS * 24 * time.Hour)
    if err = c.db.CreateSession(userIp, userInfo.Token, expiresTime); err != nil {
        http.Error(w, fmt.Sprintf("Error creating new session for username '%s'", username), http.StatusInternalServerError)
        return
    }

    /* Access Token */
    sessionInfo, numRows, err := c.db.GetSessionByUserIPAndToken(userIp, userInfo.Token)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error retrieving session for user IP '%s'", userIp), http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, fmt.Sprintf("Session not found for user IP '%s'", userIp), http.StatusUnauthorized)
        return
    }
    if sessionInfo == nil {
        http.Error(w, "Session expired", http.StatusUnauthorized)
        return
    }

    /* Cookie */
    http.SetCookie(w, &http.Cookie{
        Name: "access_token",
        Value: sessionInfo.AccessToken,
        Expires: expiresTime,
        HttpOnly: true,
    })

    /* Redirect */
    http.Redirect(w, r, serverInfo.Redirect, http.StatusSeeOther)
}