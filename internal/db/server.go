package database

import (
    "github.com/jackc/pgx/v5"
)

type ServerInfo struct {
    Name string
    Addrs []string
    Proxy string
    Redirect string
    Token string
}

func (db *Database) CreateServer(name string, addrs []string, proxy string, redirect string) error {
    query := `
        INSERT INTO servers (name, addresses, proxy_url, redirect_url)
        VALUES ($1, $2, $3, $4)`

    _, err := db.Exec(query, name, addrs, proxy, redirect)
    return err
}

func (db *Database) GetAllServerInfo() ([]ServerInfo, int, error) {
    query := `
        SELECT name, addresses, proxy_url, redirect_url, token
        FROM servers`

    var serverInfos []ServerInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        for rows.Next() {
            var serverInfo ServerInfo
            err := rows.Scan(&serverInfo.Name, &serverInfo.Addrs, &serverInfo.Proxy, &serverInfo.Redirect, &serverInfo.Token)
            if err != nil {
                continue
            }
            serverInfos = append(serverInfos, serverInfo)
        }
        return len(serverInfos), nil
    }

    numRows, err := db.Query(scanFn, query)
    return serverInfos, numRows, err
}

func (db *Database) GetServerInfoByToken(token string) (ServerInfo, int, error) {
    query := `
        SELECT name, addresses, proxy_url, redirect_url, token
        FROM servers
        WHERE token = $1`

    var serverInfo ServerInfo 
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            err := rows.Scan(&serverInfo.Name, &serverInfo.Addrs, &serverInfo.Proxy, &serverInfo.Redirect, &serverInfo.Token)
            return 1, err
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, token)
    return serverInfo, numRows, err
}

func (db *Database) GetServerInfoByID(id string) (ServerInfo, int, error) {
    query := `
        SELECT name, addresses, proxy_url, redirect_url, token
        FROM servers
        WHERE id = $1`

    var serverInfo ServerInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            err := rows.Scan(&serverInfo.Name, &serverInfo.Addrs, &serverInfo.Proxy, &serverInfo.Redirect, &serverInfo.Token)
            return 1, err
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, id)
    return serverInfo, numRows, err
}
