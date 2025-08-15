package main

import (
	"fmt"
	"net/http"
	"sippy/internal/core"
	"sync"
)

var (
	webRegistry     = core.NewRegistry()
	webCallManager  = core.NewCallManager()
	mu              sync.Mutex
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/calls", callsHandler)
	fmt.Println("Web UI running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		if username != "" {
			mu.Lock()
			webRegistry.Register(username, "webui")
			mu.Unlock()
			w.Write([]byte("<div class='container'><h1>User registered!</h1><a href='/register'>Back</a></div>"))
			return
		}
		w.Write([]byte("<div class='container'><h1>Username required!</h1><a href='/register'>Back</a></div>"))
		return
	}
	w.Write([]byte(`<html><head><link rel='stylesheet' href='/static/style.css'></head><body><div class='container'><h1>Register User</h1><form method='POST'><input name='username' placeholder='Username'><button type='submit'>Register</button></form></div></body></html>`))
}

func callsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	users := webRegistryUsers()
	calls := webCallManagerCalls()
	mu.Unlock()
	w.Write([]byte(`<html><head><link rel='stylesheet' href='/static/style.css'></head><body><div class='container'><h1>Call Monitoring</h1><h2>Registered Users</h2><ul>`))
	for _, u := range users {
		w.Write([]byte("<li>" + u.Username + "</li>"))
	}
	w.Write([]byte("</ul><h2>Active Calls</h2><ul>"))
	for _, c := range calls {
		w.Write([]byte("<li>" + c.Caller + " → " + c.Callee + "</li>"))
	}
	w.Write([]byte("</ul></div></body></html>"))
}

func webRegistryUsers() []*core.User {
	users := []*core.User{}
	for _, u := range webRegistryUsersMap() {
		users = append(users, u)
	}
	return users
}

func webRegistryUsersMap() map[string]*core.User {
	return webRegistryUsersInternal()
}

func webRegistryUsersInternal() map[string]*core.User {
	return webRegistryUsersExported(webRegistry)
}

func webRegistryUsersExported(r *core.Registry) map[string]*core.User {
	return rExported(r)
}

func rExported(r *core.Registry) map[string]*core.User {
	return r.Users()
}

func (r *core.Registry) Users() map[string]*core.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]*core.User)
	for k, v := range r.users {
		copy[k] = v
	}
	return copy
}

func webCallManagerCalls() []*core.Call {
	calls := []*core.Call{}
	for _, c := range webCallManagerCallsMap() {
		calls = append(calls, c)
	}
	return calls
}

func webCallManagerCallsMap() map[string]*core.Call {
	return webCallManagerCallsInternal()
}

func webCallManagerCallsInternal() map[string]*core.Call {
	return webCallManagerCallsExported(webCallManager)
}

func webCallManagerCallsExported(cm *core.CallManager) map[string]*core.Call {
	return cmExported(cm)
}

func cmExported(cm *core.CallManager) map[string]*core.Call {
	return cm.Calls()
}

func (cm *core.CallManager) Calls() map[string]*core.Call {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	copy := make(map[string]*core.Call)
	for k, v := range cm.calls {
		copy[k] = v
	}
	return copy
}
