package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

type Pty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

// Winsize stores the Height and Width of a terminal.
type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 // unused
	y      uint16 // unused
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h int) {
	log.Printf("window resize %dx%d", w, h)
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}

func (wp *Pty) Start(cols int, rows int) {
	var err error
	wp.Cmd = exec.Command("/bin/bash")
	wp.Pty, err = pty.Start(wp.Cmd)
	SetWinsize(wp.Pty.Fd(), cols, rows)
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
	colsString := r.URL.Query().Get("cols")
	rowsString := r.URL.Query().Get("rows")

	cols, colsParseErr := strconv.Atoi(colsString)
	rows, rowsParseErr := strconv.Atoi(rowsString)

	if colsParseErr != nil || rowsParseErr != nil {
		log.Fatalf("invalid cols/rows %s, %s\n", cols, rows)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	wp := Pty{}
	// TODO: check for errors, return 500 on fail
	wp.Start(cols, rows)

	go func() {
		// TODO: more graceful exit on socket close / process exit
		for {
			buf := make([]byte, 1024)
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
