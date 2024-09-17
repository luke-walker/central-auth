package database

func (db *Database) CreateServer(name string, addrs []string, redirect string) error {
    _, err := db.Exec("INSERT INTO servers (name, addresses, redirect_url) VALUES ($1, $2, $3)", name, addrs, redirect)
    if err != nil {
        return err
    }

    return nil
}
