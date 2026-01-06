package main

import (
	"flag"
	"fmt"
)

func main() {

	mappings := proxyMappings{}

	host := flag.String("host", "127.0.0.1", "Server host")
	port := flag.Int("port", 8080, "Server port")
	flag.Var(&mappings, "map", "Comma seperated reserve proxy mappings: /path=http://backend, /auth=http://auth")
	flag.Parse() // parse the command-line flags
	fmt.Printf("Starting reverse proxy: %s:%d with mappings %s\n", *host, *port, mappings.String())
	ProxyForever(*host, *port, mappings)
}
