package main

import (
	"time"

	"github.com/fortytw2/excalibur/process"
)

func main() {
	p, _ := process.NewProcess(&process.Config{Command: []string{"redis-server"}})

	go p.Run()

	// test out our killer instinct
	go func(proc *process.Process) {
		for {
			time.Sleep(10 * time.Second)
			proc.Kill()
		}
	}(p)

	// block until the process is over
	p.StdoutMonitor()
}
