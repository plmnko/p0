package main

import (
	"fmt"
	"github.com/cmu440/p0"
	// "time"
)

const defaultPort = 9999

func main() {
	// Initialize the server.
	server := p0.New()
	if server == nil {
		fmt.Println("New() returned a nil server. Exiting...")
		return
	}

	// Start the server and continue listening for client connections in the background.
	if err := server.Start(defaultPort); err != nil {
		fmt.Printf("MultiEchoServer could not be started: %s\n", err)
		return
	}

	// time.Sleep(10 * time.Second)
	// server.Close()

	// Block forever.
	// select {}
}
