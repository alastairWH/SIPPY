package main

import (
	"fmt"
	"net"
	"sippy/internal/sip"
	"sippy/internal/core"
	"sippy/internal/web"
)

var (
	sqliteRegistry *core.SQLiteRegistry
	registry       *core.Registry
	callManager    = core.NewCallManager()
)

func main() {
	var err error
	sqliteRegistry, err = core.NewSQLiteRegistry("sippy.db")
	if err != nil {
		panic(err)
	}
	registry = core.NewRegistryWithSQLite(sqliteRegistry)
	go runSIPServer()
	go web.StartWebUIWithRegistry(registry)
	select {} // block forever
}

func runSIPServer() {
	addr := ":5060" // SIP default port
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("SIP server listening on", addr)

	buf := make([]byte, 2048)
	for {
		n, remote, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		go handleSIPMessage(conn, remote, buf[:n])
	}
}

func handleSIPMessage(conn net.PacketConn, remote net.Addr, data []byte) {
	msg := sip.ParseSIPMessage(string(data))
	fmt.Printf("Received from %v: %s\n", remote, msg.Method)

	switch msg.Method {
	case "REGISTER":
		username := msg.Headers["To"]
		password := msg.Headers["Password"]
		fmt.Printf("Attempting to register user: %s, address: %s, password: %s\n", username, remote.String(), password)
		user := registry.GetUser(username)
		if user != nil {
			if user.Password == password {
				registry.Register(username, remote.String(), password)
				fmt.Printf("User authenticated and registered: %+v\n", user)
				// TODO: Send SIP 200 OK response
			} else {
				fmt.Printf("Authentication failed for user: %s\n", username)
				// TODO: Send SIP 403 Forbidden response
			}
		} else {
			fmt.Printf("User not found: %s\n", username)
			// TODO: Send SIP 404 Not Found response
		}
	case "INVITE":
		caller := msg.Headers["From"]
		callee := msg.Headers["To"]
		callManager.StartCall(caller, callee)
		fmt.Printf("Call started: %s -> %s\n", caller, callee)
		// TODO: Forward INVITE to callee and send SIP 200 OK response
	case "BYE":
		caller := msg.Headers["From"]
		callee := msg.Headers["To"]
		callManager.EndCall(caller, callee)
		fmt.Printf("Call ended: %s -> %s\n", caller, callee)
		// TODO: Send SIP 200 OK response
	default:
		fmt.Println("Unsupported SIP method:", msg.Method)
	}
}
