package database

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
    pool *pgxpool.Pool
}

func NewDatabase(connStr string) (*Database, error) {
    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        return nil, err
    }

    return &Database{
        pool: pool,
    }, nil
}

func (db *Database) getConnection() (*pgxpool.Conn, error) {
    conn, err := db.pool.Acquire(context.Background())
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func (db *Database) Exec(query string, params ...any) (int64, error) {
    conn, err := db.getConnection()
    if err != nil {
        return 0, err
    }
    defer conn.Release()

    commandTag, err := conn.Exec(context.Background(), query, params...)
    if err != nil {
        return 0, err
    }

    return commandTag.RowsAffected(), nil
}

func (db *Database) Query(scan func(pgx.Rows) (int, error), query string, params ...any) (int, error) {
    conn, err := db.getConnection()
    if err != nil {
        return 0, err
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(), query, params...)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    return scan(rows)
}

func (db *Database) Close() {
    db.pool.Close()
}
