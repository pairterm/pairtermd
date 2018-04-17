package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

type Pty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

func (wp *Pty) Start() {
	var err error
	wp.Cmd = exec.Command("/bin/bash")
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
}

func (wp *Pty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PtyHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	wp := Pty{}
	// TODO: check for errors, return 500 on fail
	wp.Start()

	go func() {
		// TODO: more graceful exit on socket close / process exit
		for {
			buf := make([]byte, 128)
			n, err := wp.Pty.Read(buf)
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf)

			if err != nil {
				log.Printf("Failed to send %d bytes on websocket: %s", n, err)
				return
			}
		}
	}()

	// read from the web socket, copying to the pty master
	for {
		mt, payload, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Printf("conn.ReadMessage failed: %s\n", err)
				return
			}
		}

		switch mt {
		case websocket.BinaryMessage:
			wp.Pty.Write(payload)
		case websocket.TextMessage:
			wp.Pty.Write(payload)
		default:
			log.Printf("Invalid message type %d\n", mt)
			return
		}
	}

	wp.Stop()
}
