package controllers

import (
    "net/http"
)

func errInvalidRequestMethod(w http.ResponseWriter) {
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
