package ws

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
)

type Pty struct {
	cmd *exec.Cmd
	pty *os.File
}

// winsize stores the Height and Width of a terminal.
type winsize struct {
	height uint16
	width  uint16
	x      uint16 // unused
	y      uint16 // unused
}

func (wp *Pty) Start(colsString string, rowsString string) error {
	var err error
	wp.cmd = exec.Command("/bin/bash")
	wp.pty, err = pty.Start(wp.cmd)
	if err != nil {
		return err
	}

	cols, rows, err := parseColsAndRows(colsString, rowsString)
	if err != nil {
		return err
	}
	wp.Resize(cols, rows)

	return nil
}

func (wp *Pty) Read(b []byte) (n int, err error) {
	return wp.pty.Read(b)
}

func (wp *Pty) Write(b []byte) (n int, err error) {
	return wp.pty.Write(b)
}

func (wp *Pty) Stop() {
	wp.pty.Close()
	wp.cmd.Wait()
}

func (wp *Pty) Resize(cols int, rows int) {
	log.Printf("window resize %dx%d", cols, rows)
	wins := &winsize{width: uint16(cols), height: uint16(rows)}
	syscall.Syscall(syscall.SYS_IOCTL, wp.pty.Fd(), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(wins)))
}

func parseColsAndRows(colsString string, rowsString string) (int, int, error) {
	cols, colsParseErr := strconv.Atoi(colsString)
	rows, rowsParseErr := strconv.Atoi(rowsString)

	if colsParseErr != nil || rowsParseErr != nil {
		err := fmt.Errorf("invalid cols/rows %s, %s", cols, rows)
		return -1, -1, err
	}

	return cols, rows, nil
}
