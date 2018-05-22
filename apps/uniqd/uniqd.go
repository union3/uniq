package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/union3/uniq/server"
)

var uniq *server.UNIQ

func start() {
	uniq = server.New()
	uniq.Main()
}

func stop() {
	uniq.Exit()
}

func main() {
	start()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		switch <-sig {
		case os.Interrupt:
			fallthrough
		case os.Signal(syscall.SIGTERM):
			fmt.Println("close the server")
			stop()
			os.Exit(0)
		}
	}
}
