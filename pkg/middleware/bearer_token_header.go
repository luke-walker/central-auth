package middleware

import (
    "net/http"
)

func AddBearerTokenHeader(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        accessToken, err := r.Cookie("access_token")
        if err != nil {
            http.Error(w, "Cookie 'access_token' not found", http.StatusBadRequest)
            return
        }
        r.Header.Add("Authorization", "Bearer " + accessToken.Value)

        next.ServeHTTP(w, r)
    })
}
