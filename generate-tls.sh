#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="tls"

init_dirs() {
  mkdir -p "$1"/{certs,keys}

  chmod 755 "$1"
  chmod 755 "$1/certs"
  chmod 700 "$1/keys"
}

create_ca() {
  local dir="$1"

  openssl genrsa -out "$dir/ca.key" 4096
  chmod 600 "$dir/ca.key"

  openssl req -x509 -new -nodes \
    -key "$dir/ca.key" \
    -sha256 \
    -days 3650 \
    -out "$dir/certs/ca.crt" \
    -subj "/CN=${dir##*/} Root CA"

  chmod 644 "$dir/certs/ca.crt"
}

create_key_and_csr() {
  local dir="$1"
  local name="$2"

  openssl genrsa -out "$dir/keys/$name.key" 2048
  chmod 600 "$dir/keys/$name.key"

  openssl req -new \
    -key "$dir/keys/$name.key" \
    -out "$dir/$name.csr" \
    -subj "/CN=$name"

  chmod 644 "$dir/$name.csr"
}

sign_cert() {
  local dir="$1"
  local name="$2"

  openssl x509 -req \
    -in "$dir/$name.csr" \
    -CA "$dir/certs/ca.crt" \
    -CAkey "$dir/ca.key" \
    -CAcreateserial \
    -out "$dir/certs/$name.crt" \
    -days 825 \
    -sha256 \
    -extfile "$dir/$name.ext"

  chmod 644 "$dir/certs/$name.crt"
  chmod 600 "$dir/certs/ca.srl"
}

process_proxy() {
  local dir="$1"

  echo "🔐 Processing $dir"

  if [[ ! -f "$dir/server.ext" || ! -f "$dir/client.ext" ]]; then
    echo "❌ Missing server.ext or client.ext in $dir"
    exit 1
  fi

  chmod 644 "$dir"/*.ext

  init_dirs "$dir"
  create_ca "$dir"

  create_key_and_csr "$dir" server
  create_key_and_csr "$dir" client

  sign_cert "$dir" server
  sign_cert "$dir" client
}

#######################################
# Run
#######################################

process_proxy "$BASE_DIR/reverse-proxy"
process_proxy "$BASE_DIR/forward-proxy"

echo "✅ TLS certificates generated with secure permissions"