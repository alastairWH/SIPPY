package core

import (
	"sync"
)

type User struct {
	Username string
	Address  string // network address
}

type Registry struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{users: make(map[string]*User)}
}

func (r *Registry) Register(username, address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[username] = &User{Username: username, Address: address}
}

func (r *Registry) GetUser(username string) *User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users[username]
}

func (r *Registry) Unregister(username string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, username)
}

func (r *Registry) Users() map[string]*User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]*User)
	for k, v := range r.users {
		copy[k] = v
	}
	return copy
}
