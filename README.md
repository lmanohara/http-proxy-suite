# HTTPProxySuite

**HTTPProxySuite** is a Go project demonstrating a full HTTP request flow using custom proxies and servers built with standard libraries.\
The suite includes:

1. **Forward Proxy** – receives client requests and forwards them to the reverse proxy.
2. **Reverse Proxy** – receives requests from the forward proxy and forwards them to the HTTP servers.
3. **HTTP Servers** – process requests and return responses back through the proxy chain.

This project showcases lightweight, end-to-end request handling, proxying, and server communication in Go, all without third-party frameworks.

### Highlevel Design

```mermaid
graph TD
    Client -->|HTTPS Requests| ForwardProxy[Forward Proxy]
    ForwardProxy -->|Proxied Requests| ReverseProxy[Reverse Proxy]
    ReverseProxy -->|Routed Requests| HTTPServer1[HTTP Server 1]
    ReverseProxy -->|Routed Requests| HTTPServer2[HTTP Server 2]
    HTTPServer1 -->|Responses| ReverseProxy
    HTTPServer2 -->|Responses| ReverseProxy
    ReverseProxy -->|Proxied Responses| ForwardProxy
    ForwardProxy -->|Responses| Client
```
### TLS Communication Flow

```mermaid
sequenceDiagram
    participant Client
    participant ForwardProxy as Forward Proxy <br/>Port: 6443 <br/>(TLS Server)
    participant ReverseProxy as Reverse Proxy<br/>Port: 7443<br/>(TLS Server)
    participant HTTPServer as HTTP Server<br/>Port: 8443<br/>(mTLS Server)

    Note over Client,HTTPServer: Phase 1: Client → Forward Proxy TLS

    Client->>ForwardProxy: TLS Client Hello (TLSv1.3) with CA Certificate
    ForwardProxy->>Client: Server Hello + Certificate Chain
    Note right of ForwardProxy: CN=server<br/>Issuer: forward-proxy Root CA<br/>SAN: 127.0.0.1
    ForwardProxy->>Client: Certificate Verify + Finished
    Client->>ForwardProxy: Finished
    Note over Client,ForwardProxy: TLS Connection Established

    Note over Client,HTTPServer: Phase 2: HTTP CONNECT Tunnel

    Client->>ForwardProxy: CONNECT reverse-proxy-server:7443 HTTP/1.1<br/>Host: reverse-proxy-server:7443<br/>Proxy-Connection: Keep-Alive
    ForwardProxy->>Client: HTTP/1.1 200 Connection Established
    Note over Client,ForwardProxy: CONNECT Tunnel Established

    Note over Client,HTTPServer: Phase 3: Client → Reverse Proxy TLS (through tunnel)

    Client->>ReverseProxy: TLS Client Hello (TLSv1.3)<br/>via Forward Proxy tunnel
    ReverseProxy->>Client: Server Hello + Certificate Chain<br/>via Forward Proxy tunnel
    Note right of ReverseProxy: CN=server<br/>Issuer: reverse-proxy Root CA<br/>SAN: reverse-proxy-server
    ReverseProxy->>Client: Certificate Verify + Finished<br/>via Forward Proxy tunnel
    Client->>ReverseProxy: Finished<br/>via Forward Proxy tunnel
    Note over Client,ReverseProxy: End-to-End TLS Established

    Note over Client,HTTPServer: Phase 4: HTTP Request/Response

    Client->>ReverseProxy: GET /server1 HTTP/1.1<br/>Host: reverse-proxy-server:7443<br/>Accept: */*
    
    Note over ReverseProxy,HTTPServer: Phase 5: Reverse Proxy → HTTP Server mTLS
    
    ReverseProxy->>HTTPServer: mTLS Handshake (Client Cert Required)
    HTTPServer->>ReverseProxy: Server Certificate + Client Cert Request
    ReverseProxy->>HTTPServer: Client Certificate + Finished
    HTTPServer->>ReverseProxy: Certificate Verify + Finished
    Note over ReverseProxy,HTTPServer: mTLS Connection Established
    
    ReverseProxy->>HTTPServer: GET /server1 HTTP/1.1 (over mTLS)
    HTTPServer->>ReverseProxy: HTTP/1.1 200 OK<br/>Content-Type: text/html
    
    ReverseProxy->>Client: HTTP/1.1 200 OK<br/>Content-Type: text/html<br/>Connection: close<br/><br/>Hello from Server 1

    Note over Client,HTTPServer: Certificate Chain Verification


        Note over ForwardProxy: Server Cert: CN=server<br/>CA: forward-proxy Root CA<br/>Validates: --proxy-cacert ca.crt



        Note over ReverseProxy: Server Cert: CN=server<br/>CA: reverse-proxy Root CA<br/>Validates: --cacert ca.crt



        Note over HTTPServer: mTLS Required<br/>Client + Server Certificates<br/>Mutual Authentication
```


---

## Project Structure

Only the key files are included below:

```
.
├── http-forward-proxy/
│   ├── main.go
│   └── Dockerfile
├── http-reverse-proxy/
│   ├── main.go
│   └── Dockerfile
├── http-server/
│   ├── main.go
│   ├── index_server_1.html
│   └── index_server_2.html
├── docker-compose.yml
└── README.md
```

- `http-forward-proxy/` – contains the forward proxy server main code and Dockerfile.
- `http-reverse-proxy/` – contains the reverse proxy main code and Dockerfile.
- `http-server/` – contains the HTTP server code and key index files.
- `docker-compose.yml` – orchestrates the containers and networks.

---

## Requirements

- [Docker](https://www.docker.com/get-started) >= 20.x
- [Docker Compose](https://docs.docker.com/compose/) >= 1.29.x
- Go >= 1.20 (for building the binaries)

---

## Setup & Running

### 1. Build and start all services

```bash
docker-compose up --build -d
```

### 2. Verify the services

```bash
docker ps
```

---

## Reverse Proxy Mapping

The reverse proxy mappings can be configured via the `-map` argument in the format:

```
<context_path>=<host>:<port>
```

Example in `docker-compose.yml`:

```yaml
command: ["./reverse-proxy-server", "-host", "0.0.0.0", "-port", "7090", "-map", "/server1=http-server-1:8081,/server2=http-server-2:8082"]
```

---

## Usage Examples

Server 1:

```bash
curl -v -x http://127.0.0.1:6790 http://reverse-proxy-server:7090/server1
curl -v -x https://127.0.0.1:6443 --proxy-cacert ca.crt http://reverse-proxy-server:7090/server1
```

Server 2:

```bash
curl -v -x http://127.0.0.1:6790 http://reverse-proxy-server:7090/server2
```

---

## Customizing Index Pages

- `http-server/index_server_1.html` → served by `http-server-1`
- `http-server/index_server_2.html` → served by `http-server-2`

---

## Networks

- `internal-net` – for reverse proxy and HTTP servers.
- `public-net` – exposed network for forward proxy.

---

## Stopping & Removing Containers

```bash
docker-compose down
```

---

## Notes

- Forward proxy can route requests to any reverse proxy.
- Reverse proxy mapping can be updated via `-map` argument.
- HTTP servers listen on all interfaces for container communication.

---

## License

Specify your license here (e.g., MIT, Apache 2.0).