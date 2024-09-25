package controllers

import (
    "errors"
    "net/http"
    "strings"
)

func GetBearerToken(r *http.Request) (string, error) {
    accessToken := r.Header.Get("Authorization")
    if accessToken == "" {
        return accessToken, errors.New("Missing authorization header")
    }
    if !strings.HasPrefix(accessToken, "Bearer ") {
        return accessToken, errors.New("Invalid authorization header (must be bearer)")
    }
    if len(accessToken) <= 7 {
        return accessToken, errors.New("Missing access token in authorization header")
    }
    return accessToken[7:], nil
}

func errInvalidRequestMethod(w http.ResponseWriter) {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
