package web

import (
	"fmt"
	"net/http"
	"sippy/internal/core"
	"sync"
)

var (
	registry    *core.Registry
	callManager = core.NewCallManager()
	mu          sync.Mutex
)

func StartWebUIWithRegistry(r *core.Registry) {
	registry = r
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
	http.Redirect(w, r, "/calls", http.StatusFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Printf("Web registration attempt: username=%s, password=%s\n", username, password)
		if username != "" && password != "" {
			mu.Lock()
			err := registry.Register(username, "webui", password)
			mu.Unlock()
			user := registry.GetUser(username)
			if err != nil {
				fmt.Printf("DB error on registration: %v\n", err)
				w.Write([]byte("<div class='container'><h1>Registration failed: " + err.Error() + "</h1><a href='/register'>Back</a></div>"))
				return
			}
			if user != nil {
				fmt.Printf("Web user registered: %+v\n", user)
				w.Write([]byte("<div class='container'><h1>User registered!</h1><a href='/register'>Back</a></div>"))
				return
			} else {
				fmt.Printf("Web registration failed for user: %s\n", username)
				w.Write([]byte("<div class='container'><h1>Registration failed for user: " + username + "</h1><a href='/register'>Back</a></div>"))
				return
			}
		}
		w.Write([]byte("<div class='container'><h1>Username and password required!</h1><a href='/register'>Back</a></div>"))
		return
	}
	w.Write([]byte(`<html><head><link rel='stylesheet' href='/static/style.css'></head><body><div class='container'><h1>Register User</h1><form method='POST'><input name='username' placeholder='Username'><input name='password' type='password' placeholder='Password'><button type='submit'>Register</button></form></div></body></html>`))
}

func callsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	users := registryUsers()
	fmt.Printf("Displaying registered users: %+v\n", users)
	calls := callManagerCalls()
	mu.Unlock()
	w.Write([]byte(`<html><head><link rel='stylesheet' href='/static/style.css'></head><body><div class='container'><h1>Call Monitoring</h1><h2>Registered Users</h2><ul>`))
	for _, u := range users {
		w.Write([]byte("<li><b>Username:</b> " + u.Username + "<br><b>Password:</b> " + u.Password + "<br><b>Address:</b> " + u.Address + `<br><b>Deskphone Setup:</b> SIP Server: <code>your-server-ip:5060</code>, Username: <code>` + u.Username + `</code>, Password: <code>` + u.Password + `</code></li><br>`))
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
