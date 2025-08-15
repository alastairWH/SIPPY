# SIPPY

A simple SIP server written in Go for deskphones and custom softclient communication.

## Features
- SIP REGISTER, INVITE, BYE support
- User registration and session management
- Call setup between handsets and softclient

## Getting Started
1. Install Go 1.21 or newer
2. Build and run the server:
   ```pwsh
   go run ./cmd/sippy-server/main.go
   ```

## Project Structure
- `cmd/sippy-server`: Main server entry point
- `internal/sip`: SIP protocol handling
- `internal/core`: Core logic (registration, call management)
