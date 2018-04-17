package main

import (
	"fmt"
	"net"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/pairterm/pairtermd/server"
)

var version string
var mode string

func main() {
	port := ":0"

	if mode == "dev" {
		port = ":9000"
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/pty", server.PtyHandler)
	if mode == "dev" {
		http.HandleFunc("/", server.WebpackHandler)
	} else {
		http.Handle("/", http.FileServer(rice.MustFindBox("pairtermjs/dist").HTTPBox()))
	}
	fmt.Printf("Starting pairtermd (%s) on %d...\n", version, listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
