package internal

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// StartServer starts a TCP server on the specified port.
// It listens for incoming connections and handles them concurrently.
// If the maximum number of connections is reached, new connections will be rejected.
// The server maintains a map of connected users and a channel for sending messages.
// The server runs indefinitely until an error occurs or it is manually closed.
func StartServer(port string) {
	li, errStart := net.Listen("tcp",":" +port)
	if errStart != nil {
		log.Println(portUse)
		os.Exit(0)
	}
	defer li.Close()
	fmt.Println(listenMsg + port)
	s := &server{
		listen:   li,
		messages: make(chan message),
		users:    make(map[string]net.Conn),
		mu:       sync.RWMutex{},
	}
	go s.write()
	for {
		conn, errConn := s.listen.Accept()
		if errConn != nil {
			conn.Close()
			continue
		}
		if len(s.users) > 10 {
			conn.Write([]byte(fullConn))
			for len(s.users) > 10 {
			}
		}
		go s.handler(conn)
	}
}