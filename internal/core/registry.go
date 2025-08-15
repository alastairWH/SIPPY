package core

import (
	"sync"
)

type User struct {
	Username string
	Address  string // network address
	Password string
}

type Registry struct {
	sqlite *SQLiteRegistry
	mu     sync.RWMutex
}

func NewRegistryWithSQLite(sqlite *SQLiteRegistry) *Registry {
	return &Registry{sqlite: sqlite}
}

func (r *Registry) Register(username, address, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.sqlite.Register(username, address, password)
}

func (r *Registry) GetUser(username string) *User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, err := r.sqlite.GetUser(username)
	if err != nil {
		return nil
	}
	return u
}

func (r *Registry) Unregister(username string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.sqlite.Unregister(username)
}

func (r *Registry) Users() map[string]*User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userList, err := r.sqlite.Users()
	users := make(map[string]*User)
	if err == nil {
		for _, u := range userList {
			users[u.Username] = u
		}
	}
	return users
}
