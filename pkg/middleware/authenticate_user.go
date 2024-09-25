package middleware

import (
    "net/http"

    "central-auth/internal/controllers"
    "central-auth/internal/db"
)

/* Currently only supports JSON body */
func AuthenticateUser(db *database.Database, admin bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            /* Read Headers */
            accessToken, err := controllers.GetBearerToken(r)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            /* Check Session */
            userIp := r.RemoteAddr
            sessionInfo, numRows, err := db.GetSessionByUserIPAndAccessToken(userIp, accessToken)
            if err != nil {
                http.Error(w, "Error retrieving session", http.StatusInternalServerError)
                return
            }
            if numRows == 0 {
                http.Error(w, "Could not find session", http.StatusUnauthorized)
                return
            }
            if sessionInfo == nil {
                http.Error(w, "Session has expired", http.StatusUnauthorized)
                return
            }

            if admin {
                userInfo, numRows, err := db.GetUserInfoByToken(sessionInfo.UserToken)
                if err != nil {
                    http.Error(w, "Error retrieving user info", http.StatusInternalServerError)
                    return
                }
                if numRows == 0 {
                    http.Error(w, "Could not find user", http.StatusUnauthorized)
                    return
                }
                if !userInfo.Admin {
                    http.Error(w, "User must have admin privileges", http.StatusForbidden)
                    return
                }
            }

            next.ServeHTTP(w, r)
        })
    }
}
