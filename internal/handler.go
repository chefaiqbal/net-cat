// Package internal provides the implementation of the server and its handler functions for a net-cat application.
package internal

import (
	"bufio"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

// handler is a method of the server type that handles the incoming connection.
// It performs the following steps:
// 1. Sends a welcome message to the client.
// 2. Retrieves the username from the client.
// 3. Notifies other users about the new user.
// 4. Adds the user to the list of active users.
// 5. Sends all previous messages to the client.
// 6. Handles the client's messages and broadcasts them to other users.
func (s *server) handler(conn net.Conn) {
	s.welcome(conn)
	name := s.getUserName(conn)
	s.usersNotification(conn, name)
	s.mu.Lock()
	s.users[name] = conn
	s.mu.Unlock()
	conn.Write([]byte(s.allmessages))
	s.client(conn, name)
}

// client is a method of the server type that handles the client's messages.
// It reads the messages from the client and broadcasts them to other users.
// It also adds the message to the server's message channel.
func (s *server) client(conn net.Conn, name string) {
	defer s.closeConn(conn, name)
	var text string
	buf := bufio.NewScanner(conn)
	for buf.Scan() {
		text = buf.Text()
		msg := message{
			time: time.Now().Format(dateFormat),
			user: name,
			text: text,
		}
		s.messages <- msg
	}
}

// welcome sends a welcome message to the client.
// It reads the welcome message from a file and writes it to the client's connection.
func (s *server) welcome(conn net.Conn) {
	welcom, err := os.ReadFile(welcomeMsg)
	if err != nil {
		log.Println(err)
	}
	conn.Write(welcom)
	conn.Write([]byte(nameMsg))
}

// getUserName retrieves the username from the client.
// It reads the username from the client's connection and validates it.
// If the username is invalid or already taken, it sends an error message to the client.
// It returns the valid username.
func (s *server) getUserName(conn net.Conn) string {
	var name string
	buf := bufio.NewScanner(conn)
	for buf.Scan() {
		name = buf.Text()
		if len(name) > 8 {
			conn.Write([]byte(nameIncorr + nameMsg))
		} else if !s.checkName(name) {
			conn.Write([]byte(nameIncorr + nameMsg))
		} else if _, ok := s.users[name]; ok {
			conn.Write([]byte(nameUsed + nameMsg))
		} else {
			break
		}
	}
	return name
}

// usersNotification sends a notification message to other users about the new user.
// It adds the notification message to the server's message channel.
func (s *server) usersNotification(conn net.Conn, name string) {
	msg := message{
		text: joinMsg,
		user: name,
		time: "",
	}
	s.messages <- msg
}

// closeConn closes the client's connection and removes the user from the list of active users.
// It sends a left message to other users to notify them about the user's departure.
func (s *server) closeConn(conn net.Conn, name string) {
	s.mu.Lock()
	defer conn.Close()
	defer s.mu.Unlock()
	msg := message{
		text: leftMsg,
		user: name,
		time: "",
	}
	delete(s.users, name)
	s.messages <- msg
}

// checkName checks if the given name is valid.
// It uses a regular expression to validate the name.
// It returns true if the name is valid, otherwise false.
func (s *server) checkName(name string) bool {
	pattern := regexp.MustCompile(`^[[:alpha:]]+$`)
	return pattern.MatchString(name)
}

// checkEmpty checks if the given text is empty.
// It trims the text and checks its length.
// It returns true if the text is empty, otherwise false.
func (s *server) checkEmpty(text string) bool {
	trimmedText := strings.TrimSpace(text)
	return len(trimmedText) == 0
}

// checkRune checks if the given text contains any control characters.
// It iterates over each rune in the text and checks if it is a control character.
// It returns true if the text contains control characters, otherwise false.
func (s *server) checkRune(text string) bool {
	for _, letter := range text {
		if letter < 32 {
			return true
		}
	}
	return false
}
