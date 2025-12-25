package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

var (
	CACertFilePath = "/etc/ssl/certs/ca.crt"
	CertFilePath   = "/etc/ssl/certs/server.crt"
	KeyFilePath    = "/etc/ssl/certs/server.key"
)

func ServerForever(host string, port int) {

	cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatalf("Error loading server certificate and key: %v", err)
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Error loading system cert pool: %v", err)
	}

	if caCertPEM, err := os.ReadFile(CACertFilePath); err != nil {
		log.Fatalf("Error reading CA certificate: %v", err)
	} else if ok := certPool.AppendCertsFromPEM(caCertPEM); !ok {
		log.Fatal("Error appending CA certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	address := fmt.Sprintf(":%d", port)
	listner, error := tls.Listen("tcp", address, tlsConfig)

	if error != nil {
		fmt.Print(error)
	}

	defer listner.Close()

	for {
		conn, error := listner.Accept()
		if error != nil {
			fmt.Print(error)
		}
		buffer := make([]byte, 4096) // 1kb buffer

		n, error := conn.Read(buffer)
		if error != nil {
			// conn.Close()
		}
		if n > 0 {
			byte_read := buffer[:n]
			fmt.Print(byte_read)
			responseBytes := Handle(byte_read)
			// defer conn.Close()
			conn.Write(responseBytes)
			// conn.Close()
		}
	}

}
