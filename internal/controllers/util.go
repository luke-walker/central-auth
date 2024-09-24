package controllers

import (
    "net/http"
)

func getUserIp(r *http.Request) string {
    userIp := r.Header.Get("X-Forwarded-For")
    if userIp == "" {
        userIp = r.RemoteAddr
    }
    return userIp
}

func errInvalidRequestMethod(w http.ResponseWriter) {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
