package main

type HttpRequest struct {
	Method  string
	Host    string
	Port    string
	Version string
	Headers map[string]string
}
