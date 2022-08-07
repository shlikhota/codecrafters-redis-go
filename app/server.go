package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go processConnection(c)
	}
}

func processConnection(c net.Conn) {
	scanner := bufio.NewScanner(c)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("Received from %s: %+v\n", c.RemoteAddr(), string(msg))
		if strings.ToLower(string(msg)) == "ping" {
			c.Write([]byte("PONG"))
			fmt.Printf("Sent to %s: PONG\n", c.RemoteAddr())
		}
	}
	fmt.Printf("Connection with %s has been closed!\n", c.RemoteAddr())
}
