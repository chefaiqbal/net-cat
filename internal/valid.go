package internal

import (
	"log"
)

// VaidPort validates if the given port is a valid port number.
// It checks if all characters in the port string are digits (0-9).
// If any non-digit character is found, it logs an error and exits the program.
func VaidPort(port string) {
	for _, digit := range port {
		if digit < 48 || digit > 57 {
			log.Fatalln(IncorrectPort)
		}
	}
}