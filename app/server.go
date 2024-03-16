package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
	"log"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
	// buffer := make([]byte, 1024)
	// _, err = conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading", err)
	// 	return
	// }
	// req := string(buffer)
	// firstLine := strings.Split(req, "\r\n")[0]
	// path := strings.Fields(firstLine)[1]
	// if path == "/" {
	// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// } else if strings.HasPrefix(path, "/echo") {
	// 	body := strings.SplitAfterN(path, "/", 3)
	// 	content := body[len(body)-1]
	// 	conn.Write([]byte(fmtResponseContent(content)))
	// } else if strings.HasPrefix(path, "/user-agent") {
	// 	thirdLine := strings.Split(req, "\r\n")[2]
	// 	userAgent := strings.TrimSpace(strings.Split(thirdLine, ":")[1])
	// 	conn.Write([]byte(fmtResponseContent(userAgent)))
	// } else {
	// 	conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
	// }
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading", err)
		return
	}
	req := string(buffer)
	// log.Printf("Request: %s", req)
	firstLine := strings.Split(req, "\r\n")[0]
	path := strings.Fields(firstLine)[1]
	log.Println(req)
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo") {
		body := strings.SplitAfterN(path, "/", 3)
		content := body[len(body)-1]
		conn.Write([]byte(fmtResponseContent(content)))
	} else if strings.HasPrefix(path, "/user-agent") {
		// body := strings.SplitAfterN(path, "/", 3)
		// content := body[len(body)-1]
		thirdLine := strings.Split(req, "\r\n")[2]
		userAgent := strings.TrimSpace(strings.Split(thirdLine, ":")[1])
		conn.Write([]byte(fmtResponseContent(userAgent)))
	} else {
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
	}
}

func fmtResponseContent(content string) string {
	return fmt.Sprint(fmtResponse(content) + content)
}

func fmtResponse(content string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n",
		len(content))
}