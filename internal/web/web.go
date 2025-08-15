package web

import (
	"fmt"
	"net/http"
	"sync"
	"sippy/internal/core"
)

var (
	registry     = core.NewRegistry()
	callManager  = core.NewCallManager()
	mu           sync.Mutex
)

func StartWebUI() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/calls", callsHandler)
	http.HandleFunc("/register.html", serveRegisterHTML)
	http.HandleFunc("/calls.html", serveCallsHTML)
	fmt.Println("Web UI running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveRegisterHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/register.html")
}

func serveCallsHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/calls.html")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		if username != "" {
			mu.Lock()
			registry.Register(username, "webui")
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
	users := registryUsers()
	calls := callManagerCalls()
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

func registryUsers() []*core.User {
	users := []*core.User{}
	for _, u := range registry.Users() {
		users = append(users, u)
	}
	return users
}

func callManagerCalls() []*core.Call {
	calls := []*core.Call{}
	for _, c := range callManager.Calls() {
		calls = append(calls, c)
	}
	return calls
}
