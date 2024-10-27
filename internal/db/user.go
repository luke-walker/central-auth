package database

import (
    "github.com/jackc/pgx/v5"

    "github.com/luke-walker/central-auth/internal/crypto"
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
func (db *Database) CreateUser(username string, password string, userIp string) (bool, error) { // returns if the username is available
    _, numRows, err := db.GetUserInfoByUsername(username)
    if err != nil {
        return false, err
    }
    if numRows != 0 {
        return false, nil
    }

    query := `
        INSERT INTO users (username, password, last_ip)
        VALUES ($1, $2, $3)`

    hash, err := crypto.HashPassword(password)
    if err != nil {
        return true, err
    }

    _, err = db.Exec(query, username, hash, userIp)
    return true, err
}

func (db *Database) GetUserInfoByUsername(username string) (UserInfo, int, error) {
    query := `
        SELECT id, token, username, password, last_ip, admin
        FROM users
        WHERE username = $1`

    var userInfo UserInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            err := rows.Scan(&userInfo.ID, &userInfo.Token, &userInfo.Username, &userInfo.Password, &userInfo.LastIP, &userInfo.Admin)
            return 1, err
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, username)
    return userInfo, numRows, err
}

func (db *Database) GetUserInfoByToken(token string) (UserInfo, int, error) {
    query := `
        SELECT id, token, username, password, last_ip, admin
        FROM users
        WHERE token = $1`

    var userInfo UserInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            err := rows.Scan(&userInfo.ID, &userInfo.Token, &userInfo.Username, &userInfo.Password, &userInfo.LastIP, &userInfo.Admin)
            return 1, err
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, token)
    return userInfo, numRows, err
}

func (db *Database) GetUserInfoByAccessToken(accessToken string) (UserInfo, int, error) {
    query := `
        SELECT u.token, u.username, u.last_ip
        FROM users u
        JOIN sessions s ON s.user_token = u.token
        WHERE s.access_token = $1`
    
    var userInfo UserInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            err := rows.Scan(&userInfo.Token, &userInfo.Username, &userInfo.LastIP)
            return 1, err
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, accessToken)
    return userInfo, numRows, err
}

func (db *Database) SetUserAdminStatusByToken(admin bool, token string) error {
    query := `
        UPDATE users
        SET admin = $1
        WHERE token = $2`

    _, err := db.Exec(query, admin, token)
    return err
}
