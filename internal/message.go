package internal

import (
	"fmt"
	"time"
)

// write continuously writes messages to connected users.
func (s *server) write() {
	for {
		msg := <-s.messages
		s.mu.Lock()
		text, ok := s.toString(msg)
		time := time.Now().Format(dateFormat)
		var datename string
		for name, userAddr := range s.users {
			datename = fmt.Sprintf("\n[%s][%s]:", time, name)
			if msg.user == name {
				userAddr.Write([]byte(datename[1:]))
				continue
			}
			if ok {
				userAddr.Write([]byte(text))
				userAddr.Write([]byte(datename))
			}
		}
		s.mu.Unlock()
	}
}

// toString converts a message to a string format.
// It returns the converted string and a boolean indicating if the message is valid.
func (s *server) toString(msg message) (string, bool) {
	if s.checkEmpty(msg.text) || s.checkRune(msg.text) {
		return fmt.Sprintf("\n%s%s", msg.user, msg.text), false
	}
	if msg.time == "" {
		return fmt.Sprintf("\n%s%s", msg.user, msg.text), true
	}
	text := fmt.Sprintf("\n[%s][%s]:%s", msg.time, msg.user, msg.text)
	s.saveMessage(text[1:] + "\n")
	return text, true
}

// saveMessage saves a message to the server's message history.
func (s *server) saveMessage(msg string) {
	s.allmessages += msg
}
