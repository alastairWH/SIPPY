package sip

import (
	"strings"
)

// SIPMessage represents a parsed SIP message
type SIPMessage struct {
	Method   string
	Headers  map[string]string
	Body     string
}

// ParseSIPMessage parses raw SIP message data into a SIPMessage struct
func ParseSIPMessage(data string) *SIPMessage {
	lines := strings.Split(data, "\r\n")
	msg := &SIPMessage{Headers: make(map[string]string)}
	if len(lines) == 0 {
		return msg
	}
	// First line: Method (e.g., REGISTER, INVITE)
	parts := strings.Split(lines[0], " ")
	if len(parts) > 0 {
		msg.Method = parts[0]
	}
	// Headers
	bodyStart := false
	for _, line := range lines[1:] {
		if line == "" {
			bodyStart = true
			continue
		}
		if bodyStart {
			msg.Body += line + "\n"
		} else {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) == 2 {
				msg.Headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}
	return msg
}
