package main

import (
	"crypto/tls"
	"fmt"
	"os"
)

var (
	ServerCertFilePath = os.Getenv("TLS_SERVER_CERT")
	ServerKeyFilePath  = os.Getenv("TLS_SERVER_KEY")
)

func ProxyForever(host string, port int, mappings proxyMappings) {
	cert, err := tls.LoadX509KeyPair(ServerCertFilePath, ServerKeyFilePath)

	if err != nil {
		fmt.Println("Error loading server certificate and key:", err)
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert, // Don't require client certificates
	}

	address := fmt.Sprintf(":%d", port)

	listener, error := tls.Listen("tcp", address, tlsConfig)

	if error != nil {
		fmt.Println(error)
	}

	defer listener.Close()

	for {
		conn, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		buffer := make([]byte, 4096)
		n, error := conn.Read(buffer)
		if error != nil {
			fmt.Println(error)
		}

		if n > 0 {
			byte_read := buffer[:n]
			out := Handle(byte_read, mappings)
			fmt.Print(byte_read)
			conn.Write(out)
			// // conn.Close()
		}
	}

}
