package main

import (
	"log"
	"os"

	"net-cat/internal"
)

// main is the entry point of the program.
// It initializes the port variable with the default port value,
// reads the command-line arguments, and starts the server.
func main() {
	port := internal.DefPort
	args := os.Args[1:]
	if len(args) > 1 {
		log.Fatalln(internal.IncorrectPort)
	}
	if len(args) == 1 {
		port = args[0]
		internal.VaidPort(port)
	}
	internal.StartServer(port)
}