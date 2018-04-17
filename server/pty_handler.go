package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pairterm/pairtermd/ws"
)

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

	wp := ws.Pty{}
	err = wp.Start(r.URL.Query().Get("cols"), r.URL.Query().Get("rows"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error starting terminal: %+v\n", err)))
		return
	}

	// read from the pty, copying to the websocket
	go streamFromPty(wp, conn)
	// read from the web socket, copying to the pty
	streamToPty(wp, conn)

	wp.Stop()
}

func streamFromPty(wp ws.Pty, conn *websocket.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := wp.Read(buf)
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
}

func streamToPty(wp ws.Pty, conn *websocket.Conn) {
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
			wp.Write(payload)
		case websocket.TextMessage:
			wp.Write(payload)
		default:
			log.Printf("Invalid message type %d\n", mt)
			return
		}
	}
}
