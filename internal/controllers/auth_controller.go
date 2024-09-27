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

    serverToken := chi.URLParam(r, "serverToken")
    if serverToken == "" {
        http.Error(w, "Missing server token as URL parameter", http.StatusBadRequest)
        return
    }

    /* Verify Server */
    _, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, "Could not find matching server token", http.StatusUnauthorized)
        return
    }

    /* Send login page */
    loginTemplate, err := template.ParseFiles(fmt.Sprintf("%s/internal/templates/login.html.tmpl", os.Getenv("PROJECT_PATH")))
    if err != nil {
        http.Error(w, "Error parsing login.html template", http.StatusInternalServerError)
        return
    }
    loginTemplate.Execute(w, struct{
        ServerToken string
    }{
        ServerToken: serverToken,
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

    serverToken := chi.URLParam(r, "serverToken")
    if serverToken == "" {
        http.Error(w, "Missing server token as URL parameter", http.StatusBadRequest)
        return
    }

    /* Verify Server */
    serverInfo, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, "Could not find matching server token", http.StatusUnauthorized)
        return
    }

    /* Verify User */
    userIp := r.RemoteAddr
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
        Path: "/",
        Expires: expiresTime,
        HttpOnly: true,
        SameSite: 3,
    })

    /* Redirect */
    http.Redirect(w, r, serverInfo.Redirect, http.StatusSeeOther)
}

func (c *AuthController) AttemptUserSignUp(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errInvalidRequestMethod(w)
        return
    }
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing URL-encoded form", http.StatusInternalServerError)
        return
    }

    serverToken := chi.URLParam(r, "serverToken")
    if serverToken == "" {
        http.Error(w, "Missing server token as URL parameter", http.StatusBadRequest)
        return
    }

    /* Verify Server */
    _, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, "Could not find matching server token", http.StatusUnauthorized)
        return
    }

    /* Create User */
    userIp := r.RemoteAddr
    username := r.PostForm["username"][0]
    password := r.PostForm["password"][0]
    available, err := c.db.CreateUser(username, password, userIp)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }
    if !available {
        http.Error(w, "Username already taken", http.StatusUnauthorized)
        return
    }

    /* Login User */
    c.AttemptUserLogin(w, r)
}

func (c *AuthController) AttemptUserLogout(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errInvalidRequestMethod(w)
        return
    }

    /*
    serverToken := chi.URLParam(r, "serverToken")
    if serverToken == "" {
        http.Error(w, "Missing server token as URL parameter", http.StatusBadRequest)
        return
    }
    */

    /* Verify Server */
    /*
    _, numRows, err := c.db.GetServerInfoByToken(serverToken)
    if err != nil {
        http.Error(w, "Error retrieving server information", http.StatusInternalServerError)
        return
    }
    if numRows == 0 {
        http.Error(w, "Could not find matching server token", http.StatusUnauthorized)
        return
    }
    */

    /* Access Token */
    accessToken, err := GetBearerToken(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    /* Destroy Session */
    if _, err = c.db.DestroySessionByAccessToken(accessToken); err != nil {
        http.Error(w, "Error destroying session", http.StatusInternalServerError)
        return
    }

    /* Delete Cookie */
    http.SetCookie(w, &http.Cookie{
        Name: "access_token",
        Value: "",
        Path: "/",
        MaxAge: -1,
    })
}
