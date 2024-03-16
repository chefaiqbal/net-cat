package internal

import (
	"net"
	"sync"
)

// DefPort is the default port number used by the application.
const (
	DefPort       = "8989"
	portUse       = "This port is already in use"
	IncorrectPort = "[USAGE]: ./TCPChat $port"
	dateFormat    = "2006-01-02 15:04:05"
	joinMsg       = " has joined our chat..."
	leftMsg       = " has left our chat..."
	nameUsed      = "There's already a user with that name.\n"
	nameIncorr    = "Use only Latin letters.\n"
	nameVeryLong  = "Your name must not exceed 10 Latin characters.\n"
	nameMsg       = "[ENTER YOUR NAME]:"
	listenMsg     = "Listening on the port :"
	welcomeMsg    = "static/welcome.txt"
	fullConn      = "The server's full. Do you want to wait for someone to come out?\n"
	exitServer    = "Sorry. The server has been suspended"
)

// message represents a chat message.
type message struct {
	time string // The timestamp of the message.
	user string // The username of the sender.
	text string // The content of the message.
}

// server represents the chat server.
type server struct {
	listen      net.Listener         // The network listener for accepting incoming connections.
	messages    chan message         // The channel for sending and receiving chat messages.
	users       map[string]net.Conn  // The map of connected users, with their usernames as keys and network connections as values.
	allmessages string               // The concatenated string of all chat messages.
	mu          sync.RWMutex         // The mutex for synchronizing access to the server's state.
}