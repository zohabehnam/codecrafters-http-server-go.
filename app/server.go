package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
	"flag"
	"bytes"
	"errors"
	"io"
)

var dirFlag *string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	dirFlag = flag.String("directory", "", "Specify directory")
	flag.Parse()

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "localhost:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	defer l.Close()
	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		req := make([]byte, 1024)
		_, err := conn.Read(req)
		if errors.Is(err, io.EOF) {
			fmt.Println("Client broke the connection")
			break
		} else if err != nil {
			fmt.Println("Error reading the request")
			break
		}
		content := strings.Split(string(req), "\r\n")
		if len(content) == 0 {
			continue
		}
		requestLine := content[0]
		data := strings.Split(requestLine, " ")
		if len(data) < 2 {
			continue
		}
		method := data[0]
		path := data[1]
		if strings.HasPrefix(path, "/files") {
			if method == "POST" {
				handlePostFiles(path, content[len(content)-1], conn)
			} else {
				handleGetFiles(path, conn)
			}
		} else if strings.HasPrefix(path, "/user-agent") {
			handleUserAgent(content, conn)
		} else if strings.HasPrefix(path, "/echo/") {
			handleEcho(path, conn)
		} else if path != "/" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
	}
}

func handleEcho(path string, conn net.Conn) {
	body := strings.SplitAfterN(path, "/", 3)
	content := body[len(body)-1]
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
}

func fmtResponseContent(content string) string {
	return fmt.Sprint(fmtResponse(content, "text/plain") + content)
}

func handlePostFiles(path string, content string, conn net.Conn) {
	dir := *dirFlag
	arrPath := strings.SplitAfterN(path, "/", 3)
	finalPath := arrPath[len(arrPath)-1]
	err := os.WriteFile(dir+"/"+finalPath, bytes.Trim([]byte(content), "\x00"), 0666)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	conn.Write([]byte("HTTP/1.1 201 OK\r\nContent-Length: 0\r\n\r\n"))
}


func fmtResponse(content string, contentType string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n",
		contentType,
		len(content))
}

func handleGetFiles(path string, conn net.Conn) {
	dir := *dirFlag
	body := strings.SplitAfterN(path, "/", 3)
	pathToFile := body[len(body)-1]
	file, err := os.ReadFile(dir + "/" + pathToFile)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), file)))
}
func handleUserAgent(request []string, conn net.Conn) {
	content := ""
	for _, line := range request {
		if strings.Contains(line, "User-Agent") {
			fmt.Println(line)
			content = strings.Split(line, ": ")[1]
			break
		}
	}
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
}