package main

import (
	"fmt"
	"net"
	"sippy/internal/core"
	"sippy/internal/sip"
)

var registry = core.NewRegistry()
var callManager = core.NewCallManager()

func main() {
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
		registry.Register(username, remote.String())
		fmt.Printf("Registered user: %s at %s\n", username, remote.String())
		// TODO: Send SIP 200 OK response
	case "INVITE":
		caller := msg.Headers["From"]
		callee := msg.Headers["To"]
		callManager.StartCall(caller, callee)
		fmt.Printf("Call started: %s -> %s\n", caller, callee)
		// TODO: Send SIP 200 OK response and notify callee
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
