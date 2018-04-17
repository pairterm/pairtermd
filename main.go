package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/pairterm/pairtermd/server"
)

var version string

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/pty", server.PtyHandler)
	http.HandleFunc("/", server.WebpackHandler)
	fmt.Printf("Starting pairtermd on %d...\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
