package main

import (
	"crypto/tls"
	"fmt"
	"os"
)

var (
	CertFilePath = os.Getenv("TLS_SERVER_CERT")
	KeyFilePath  = os.Getenv("TLS_SERVER_KEY")
)

func ProxyForever(host string, port int) {

	cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)

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

	// defer listener.Close()

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
			err := HandleConnect(conn, byte_read)
			if err != nil {
				fmt.Println(err)
			}
			// conn.Write(out)
			// // conn.Close()
		}
	}

}
