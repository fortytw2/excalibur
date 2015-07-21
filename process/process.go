package process

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

// message is used to wrap sending things like
type message int

const (
	// killCmd tells a process to end - iota doesn't work due to `chan int` not
	// playing well with 0 values
	killCmd message = 1
)

// ByteChannelWriter provides a newline buffered io.Writer that writes each line
// to the given channel - it is NOT ok to create a ByteChannelWriter using
// anything other than NewByteChannelWriter()
type ByteChannelWriter struct {
	buf *bytes.Buffer
	c   chan []byte
}

// NewByteChannelWriter creates a new ByteChannelWriter on the given chan []byte
func NewByteChannelWriter(ch chan []byte) *ByteChannelWriter {
	cw := &ByteChannelWriter{c: ch, buf: new(bytes.Buffer)}
	go cw.flush()
	return cw
}

func (cw *ByteChannelWriter) flush() {
	for {
		time.Sleep(10 * time.Millisecond)
		line, err := cw.buf.ReadBytes(byte(0x0A))
		if err != nil {
			cw.buf.Write(line)
		}

		cw.c <- line
	}

}

func (cw *ByteChannelWriter) Write(p []byte) (n int, err error) {
	return cw.buf.Write(p)
}

// A Process is an instance of a program
type Process struct {
	Command []string
	Stdin   chan []byte
	Stdout  chan []byte
	Stderr  chan []byte

	execCmd *exec.Cmd
	control chan message
}

// NewProcess creates a new process from the given config
func NewProcess(conf *Config) (*Process, error) {
	p := &Process{
		Command: conf.Command,
		Stdin:   make(chan []byte),
		Stdout:  make(chan []byte),
		Stderr:  make(chan []byte),

		execCmd: new(exec.Cmd),
		control: make(chan message),
	}
	if len(p.Command) >= 2 {
		p.execCmd = exec.Command(p.Command[0], p.Command[1:]...)
	} else {
		p.execCmd = exec.Command(p.Command[0])
	}

	p.execCmd.Stdout = NewByteChannelWriter(p.Stdout)
	p.execCmd.Stderr = NewByteChannelWriter(p.Stderr)

	return p, nil
}

// Run launches the process along with its control channel
func (p *Process) Run() {
	go p.controlMonitor()
	p.execCmd.Run()
}

// StdoutMonitor just prints the process Std{out,err} to os.Stdout
// helpful for debugging/testing individual processes
func (p *Process) StdoutMonitor() {
	for {
		select {
		case msg, ok := <-p.Stdout:
			if !ok {
				return
			}
			fmt.Print(string(msg))
		}
	}
}

// controlMonitor handles writings on the control channel - mostly involving
// sending various signals to a running process
func (p *Process) controlMonitor() {
	for {
		select {
		case msg, ok := <-p.control:
			if !ok {
				return
			}
			if msg == killCmd && p.execCmd.Process != nil {
				p.execCmd.Process.Kill()

				close(p.Stdout)
				close(p.Stderr)
				close(p.control)
			}
		}
	}
}

// Kill sends a kill command to the control channel
func (p *Process) Kill() {
	p.control <- killCmd
}
