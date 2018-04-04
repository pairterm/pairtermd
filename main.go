package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting pairtermd on %d...\n", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, nil))
}
