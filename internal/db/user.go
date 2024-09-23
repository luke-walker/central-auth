package database

import (
    "github.com/jackc/pgx/v5"

    "central-auth/internal/crypto"
)

type UserInfo struct {
    ID string
    Token string
    Username string
    Password string
    LastIP string
    Admin bool
}

/* TODO: Insert user IP as last_ip */
func (db *Database) CreateUser(username string, password string) error {
    query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)`

    hash, err := crypto.HashPassword(password)
    if err != nil {
        return err
    }

    _, err = db.Exec(query, username, hash)
    return err
}

func (db *Database) GetUserInfoByUsername(username string) (UserInfo, int, error) {
    query := `
        SELECT id, token, username, password, last_ip, admin
        FROM users
        WHERE username = $1`

    var userInfo UserInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            rows.Scan(&userInfo.ID, &userInfo.Token, &userInfo.Username, &userInfo.Password, &userInfo.LastIP, &userInfo.Admin)
            return 1, nil
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, username)
    return userInfo, numRows, err
}
