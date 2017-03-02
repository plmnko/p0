package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	conn, err := net.Dial("tcp", defaultHost+":"+fmt.Sprintf("%d", defaultPort))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		go read_data(conn)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text+"\n")

	}
}

func read_data(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
		if err != nil {
			break
		}
	}

}
