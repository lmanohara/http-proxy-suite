package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	CACertFilePath = os.Getenv("TLS_CA")
	CertFilePath   = os.Getenv("TLS_SERVER_CERT")
	KeyFilePath    = os.Getenv("TLS_SERVER_KEY")
)

func Handle(buff []byte, mapping proxyMappings) []byte {

	cert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		fmt.Println("Error loading client certificate and key:", err)
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
		RootCAs:      certPool,
	}

	HttpRequest, error := parsedRequest(buff)
	if error != nil {
		fmt.Println(error)
	}
	path := HttpRequest.Path
	parsedUrl, error := url.Parse(path)
	if error != nil {
		fmt.Println(error)
	}

	contextPath := parsedUrl.Path

	sourceAddress, key := mapping[contextPath]

	if !key {
		var buf bytes.Buffer
		writeResponseLine(&buf, 502)
		return buf.Bytes()
	}

	fmt.Println("source address: ", sourceAddress)
	// establish connection to the http server host and port
	// conn, error := net.Dial("tcp", sourceAddress)
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 15 * time.Second}, "tcp", sourceAddress, tlsConfig)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer conn.Close()

	forwardRequest("/", HttpRequest, conn)

	buffer := make([]byte, 4096)
	n, error := conn.Read(buffer)
	if error != nil {
		fmt.Println(error)
	}

	return buffer[:n]
}

func writeResponseLine(buff *bytes.Buffer, statusCode int) {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, getStatusText(statusCode))
	buff.WriteString(statusLine)
}

func getStatusText(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 502:
		return "Bad Gateway"
	default:
		return ""
	}
}

func forwardRequest(contextPath string, HttpRequest HttpRequest, conn net.Conn) {
	var buf bytes.Buffer
	responseLine := fmt.Sprintf("GET %s %s\r\n", contextPath, HttpRequest.Version)
	buf.WriteString(responseLine)

	for k, v := range HttpRequest.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	requestBytes := buf.Bytes()

	fmt.Println("Response as string:\n", string(requestBytes))

	// forward request headers to the http server
	conn.Write(requestBytes)
}

func parsedRequest(inputStream []byte) (HttpRequest, error) {
	data := string(inputStream)
	lines := strings.Split(data, "\r\n")
	if len(lines) == 0 {
		return HttpRequest{}, errors.New("empty request")
	}

	requestLine := lines[0]
	requestLineParts := strings.SplitAfterN(requestLine, " ", 3)
	command := strings.TrimSpace(requestLineParts[0])
	path := strings.TrimSpace(requestLineParts[1])
	version := strings.TrimSpace(requestLineParts[2])

	fmt.Println("Request line: ", requestLine)
	headers := make(map[string]string)

	for _, line := range lines[1:] {
		if line == "" {
			break
		}

		keyValue := strings.SplitN(line, ":", 2)

		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			headers[key] = value
		}
	}

	for key, val := range headers {
		fmt.Printf("%s: %s\n", key, val)
	}

	req := HttpRequest{
		Method:  command,
		Path:    path,
		Version: version,
		Headers: headers,
	}

	return req, nil
}
