terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "4.4.0"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Create the network if it does not already exist.

resource "docker_network" "internal_net" {
  name = "internal-net"
}

resource "docker_image" "http_server_1" {
  name = "http-proxy-suite-http-server-1:latest"
  
  keep_locally = true
}

resource "docker_container" "http_server_1" {
  name  = "tf-http-server-1"
  image = docker_image.http_server_1.image_id

  command = ["./server", "-host", "0.0.0.0", "-port", "8443"]

  restart = "unless-stopped"

  ports {
    internal = 8443
    external = 8443
  }
  networks_advanced {
    name = docker_network.internal_net.name
  }

  env = [
    "TLS_SERVER_CERT=/etc/ssl/certs/server.crt",
    "TLS_SERVER_KEY=/etc/ssl/keys/server.key",
    "TLS_CA=/etc/ssl/certs/ca.crt"
  ]

  mounts {
    target    = "/root/index.html"
    source    = abspath("${path.module}/../http-server/index_server_1.html")
    type      = "bind"
    read_only = true
  }

  # ./tls-configs/web-server/certs/server.crt:/etc/ssl/certs/server.crt:ro

  mounts {
    target    = "/etc/ssl/certs/server.crt"
    source    = abspath("${path.module}/../tls-configs/web-server/certs/server.crt")
    type      = "bind"
    read_only = true
  }

  # ./tls-configs/web-server/certs/ca.crt:/etc/ssl/certs/ca.crt:ro

  mounts {
    target    = "/etc/ssl/certs/ca.crt"
    source    = abspath("${path.module}/../tls-configs/web-server/certs/ca.crt")
    type      = "bind"
    read_only = true
  }

  # ./tls-configs/web-server/keys/server.key:/etc/ssl/keys/server.key:ro

  mounts {
    target    = "/etc/ssl/keys/server.key"
    source    = abspath("${path.module}/../tls-configs/web-server/keys/server.key")
    type      = "bind"
    read_only = true
  }
}
