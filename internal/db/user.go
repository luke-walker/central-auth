package database

import (
    "github.com/jackc/pgx/v5"

    "central-auth/internal/crypto"
)

type UserInfo struct {
    Username string
    Password string
    LastIP string
}

/* TODO: Insert user IP */
func (db *Database) CreateUser(username string, password string) error {
    hash, err := crypto.HashPassword(password)
    if err != nil {
        return err
    }

    query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)`

    _, err = db.Exec(query, username, hash)
    return err
}

func (db *Database) GetUserInfoByUsername(username string) (UserInfo, int, error) {
    query := `
        SELECT username, password, last_ip
        FROM users
        WHERE username = $1`

    var userInfo UserInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            rows.Scan(&userInfo.Username, &userInfo.Password, &userInfo.LastIP)
            return 1, nil
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, username)
    return userInfo, numRows, err
}
