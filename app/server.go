package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type request struct {
	message []string
}

func (r *request) parse(scanner *bufio.Scanner, size int) {
	for i := 0; i < size; i++ {
		if ok := scanner.Scan(); !ok {
			return
		}

		dataInfo := scanner.Text()
		dataType, strSize := dataInfo[0], dataInfo[1:]
		if dataType != '$' {
			fmt.Printf("Wrong type of data (%c), only string ($) is allowed\n", dataType)
			return
		}

		if ok := scanner.Scan(); !ok {
			return
		}

		data := scanner.Text()
		size, err := strconv.Atoi(string(strSize))
		if err != nil || len(data) != size {
			fmt.Printf("Error parsing data (len=%d): %s\n", size, err)
			return
		}
		r.message = append(r.message, data)
	}
	// fmt.Printf("Received: %#v\n", r.message)
}

func (r *request) command() (command string) {
	if len(r.message) < 1 {
		return
	}
	return strings.ToLower(r.message[0])
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	storage := NewStorage()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go processConnection(c, storage)
	}
}

func processConnection(c net.Conn, s *storage) {
	scanner := bufio.NewScanner(c)
	scanner.Split(bufio.ScanLines)
	for {
		ok := scanner.Scan()
		msg := scanner.Text()
		if !ok || len(msg) < 2 {
			if scanner.Err() == nil {
				fmt.Printf("Connection with %s has been closed!\n", c.RemoteAddr())
				return
			}
			fmt.Printf("Error while reading: %s\n", scanner.Err())
			continue
		}

		msgType, msgSize := msg[0], msg[1:]
		switch msgType {
		// Parse array
		case '*':
			size, err := strconv.Atoi(string(msgSize))
			if err != nil {
				fmt.Printf("Error while reading size of array: %s\n", err)
			}
			r := request{}
			r.parse(scanner, size)
			c.Write(proccessRequest(r, s))
		default:
			c.Write(errorResponse("unknown prefix '%+v'", msgType))
		}
	}
}

func proccessRequest(req request, s *storage) (response []byte) {
	cmd := req.command()
	switch cmd {
	case "ping":
		response = []byte("+PONG\r\n")
	case "echo":
		response = buildBulkString(req.message[1])
	case "set":
		var args []string
		if len(req.message) < 3 {
			return errorResponse("Wrong arguments for SET command")
		} else if len(req.message) > 3 {
			args = req.message[3:]
		}
		s.Set(req.message[1], req.message[2], args)
		response = buildSimpleString("OK")
	case "get":
		key := req.message[1]
		if value, err := s.Get(key); err == nil {
			return buildBulkString(value)
		} else {
			return buildNullBulkString()
		}
	default:
		return errorResponse("Unknown command: %s", cmd)
	}
	return
}

func errorResponse(message string, params ...interface{}) []byte {
	return []byte("-ERR " + fmt.Sprintf(message, params...) + "\r\n")
}

func buildNullBulkString() (resp []byte) {
	return []byte("$-1\r\n")
}

func buildBulkString(s string) (resp []byte) {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func buildSimpleString(s string) (resp []byte) {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}
