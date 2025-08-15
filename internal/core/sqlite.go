package core

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type SQLiteRegistry struct {
	db *sql.DB
}

func NewSQLiteRegistry(dbPath string) (*SQLiteRegistry, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, address TEXT, password TEXT)`); err != nil {
		return nil, err
	}
	return &SQLiteRegistry{db: db}, nil
}

func (r *SQLiteRegistry) Register(username, address, password string) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO users (username, address, password) VALUES (?, ?, ?)`, username, address, password)
	return err
}

func (r *SQLiteRegistry) GetUser(username string) (*User, error) {
	row := r.db.QueryRow(`SELECT username, address, password FROM users WHERE username = ?`, username)
	var u User
	if err := row.Scan(&u.Username, &u.Address, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLiteRegistry) Unregister(username string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE username = ?`, username)
	return err
}

func (r *SQLiteRegistry) Users() ([]*User, error) {
	rows, err := r.db.Query(`SELECT username, address, password FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Address, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
