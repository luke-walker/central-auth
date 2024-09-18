package database

import (
    "github.com/jackc/pgx/v5"
)

type ServerInfo struct {
    Name string
    Addrs []string
    Redirect string
}

func (db *Database) CreateServer(name string, addrs []string, redirect string) error {
    query := `
        INSERT INTO servers (name, addresses, redirect_url)
        VALUES ($1, $2, $3)`

    var err error
    if redirect == "" {
        _, err = db.Exec(query, name, addrs, nil)
    } else {
        _, err = db.Exec(query, name, addrs, redirect)
    }
    return err
}

func (db *Database) GetAllServerInfo() ([]ServerInfo, int, error) {
    query := `
        SELECT name, addresses, redirect_url
        FROM servers`

    var serverInfos []ServerInfo
    scanFn := func(rows pgx.Rows) (int, error) {
        for rows.Next() {
            var serverInfo ServerInfo
            rows.Scan(&serverInfo.Name, &serverInfo.Addrs, &serverInfo.Redirect)
            serverInfos = append(serverInfos, serverInfo)
        }
        return len(serverInfos), nil
    }

    numRows, err := db.Query(scanFn, query)
    return serverInfos, numRows, err
}

func (db *Database) GetServerInfoByName(name string) (ServerInfo, int, error) {
    query := `
        SELECT name, addresses, redirect_url
        FROM servers
        WHERE name = $1`

    var serverInfo ServerInfo 
    scanFn := func(rows pgx.Rows) (int, error) {
        if rows.Next() {
            rows.Scan(&serverInfo.Name, &serverInfo.Addrs, &serverInfo.Redirect)
            return 1, nil
        }
        return 0, nil
    }

    numRows, err := db.Query(scanFn, query, name)
    return serverInfo, numRows, err
}
