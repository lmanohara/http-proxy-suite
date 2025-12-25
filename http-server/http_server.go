package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

var (
	CACertFilePath = os.Getenv("TLS_CA")
	CertFilePath   = os.Getenv("TLS_SERVER_CERT")
	KeyFilePath    = os.Getenv("TLS_SERVER_KEY")
)

func ServerForever(host string, port int) {

	cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		fmt.Println("Error loading server certificate and key:", err)
		panic(err)
	}

	certPool, err := x509.SystemCertPool()

	if err != nil {
		fmt.Println("Error loading system cert pool:", err)
		panic(err)
	}

	if caCertPEM, err := os.ReadFile(CACertFilePath); err != nil {
		fmt.Println("Error reading CA certificate:", err)
		panic(err)
	} else if ok := certPool.AppendCertsFromPEM(caCertPEM); !ok {
		fmt.Println("Error appending CA certificate to pool")
		panic(err)
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
