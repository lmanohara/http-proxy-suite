package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func HandleConnect(clientConn net.Conn, inputStream []byte) error {
	data := string(inputStream)

	lines := strings.Split(data, "\r\n")
	if len(lines) == 0 {
		return errors.New("empty request")
	}

	requestLine := lines[0]
	requestLineParts := strings.SplitAfterN(requestLine, " ", 3)
	command := strings.TrimSpace(requestLineParts[0])

	var host, port string
	if command == "CONNECT" {
		httpsUrl := strings.TrimSpace(requestLineParts[1])
		httpsUrlParts := strings.Split(httpsUrl, ":")
		host = httpsUrlParts[0]
		port = httpsUrlParts[1]
	}

	networkAddress := net.JoinHostPort(host, port)
	fmt.Println("Connecting to ", networkAddress)
	upstreamConn, error := net.Dial("tcp", networkAddress)

	if error != nil {
		fmt.Println(error)
		return error
	}

	// send connection established response to the client
	response := "HTTP/1.1 200 Connection Established\r\n\r\n"
	clientConn.Write([]byte(response))

	go func() {
		_, err := io.Copy(upstreamConn, clientConn)
		log.Println("client → upstream:", err)
	}()

	go func() {
		_, err := io.Copy(clientConn, upstreamConn)
		log.Println("upstream → client:", err)
	}()

	return nil
}
