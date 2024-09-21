package database

import (
    "errors"
    "time"

    "github.com/jackc/pgx/v5"
)

type SessionInfo struct {
    UserIP string
    UserToken string
    AccessToken string
    Expires time.Time
}

func (db *Database) CreateSession(userIp string, userToken string, expires time.Time) error {
    query := `
        INSERT INTO sessions (user_ip, user_token, expires)
        VALUES ($1, $2, $3)`

    _, err := db.Exec(query, userIp, userToken, expires)
    return err
}

func (db *Database) CheckSession(userIp string, accessToken string) (bool, error) {
    query := `
        SELECT expires
        FROM sessions
        WHERE user_ip = $1 AND access_token = $2`

    var sessionInfo SessionInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            rows.Scan(&sessionInfo.Expires)
            return 1, nil
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, userIp, accessToken)
    if err != nil || numRows == 0 {
        return false, err
    }
    if time.Now().After(sessionInfo.Expires) {
        return false, errors.New("Expired access token")
    }

    return true, nil
}

func (db *Database) GetSessionByUserIPAndToken(userIp string, userToken string) (*SessionInfo, int, error) { // SessionInfo nil if expired/error
    query := `
        SELECT user_ip, user_token, access_token, expires
        FROM sessions
        WHERE user_ip = $1 AND user_token = $2`

    var sessionInfo SessionInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            rows.Scan(&sessionInfo.UserIP, &sessionInfo.UserToken, &sessionInfo.AccessToken, &sessionInfo.Expires)
            return 1, nil
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, userIp, userToken)
    if err != nil || numRows == 0 {
        return nil, numRows, err
    }

    /* Check Expires */
    if time.Now().After(sessionInfo.Expires) {
        return nil, 0, nil
    }

    return &sessionInfo, numRows, nil
}

func (db *Database) CheckSessionExists(userIp string, userToken string) (bool, error) {
    if sessionInfo, numRows, err := db.GetSessionByUserIPAndToken(userIp, userToken); sessionInfo == nil || numRows == 0 || err != nil {
        return false, err
    }
    return true, nil
}

func (db *Database) DestroySession(userIp string, userToken string) (bool, error) {
    query := `
        DELETE FROM sessions
        WHERE user_ip = $1 AND user_token = $2`

    numRows, err := db.Exec(query, userIp, userToken)
    if err != nil || numRows == 0 {
        return false, err
    }
    return true, nil
}
