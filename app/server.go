package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading", err)
		return
	}
	firstLine := strings.Split(string(buffer), "\r\n")[0]
	path := strings.Fields(firstLine)[1]
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo") {
		body := strings.SplitAfterN(path, "/", 3)
		content := body[len(body)-1]
		conn.Write([]byte(fmtResponseContent(content)))
	} else {
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
	}
}

func fmtResponseContent(content string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(content),
		content)
}