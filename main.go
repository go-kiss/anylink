package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kiss/anylink/base"
	"github.com/go-kiss/anylink/handler"
)

// 程序版本
var COMMIT_ID string

func main() {
	base.CommitId = COMMIT_ID

	base.Start()
	handler.Start()
	signalWatch()
}

func signalWatch() {
	base.Info("Server pid: ", os.Getpid())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
	for {
		sig := <-sigs
		base.Info("Get signal:", sig)
		switch sig {
		case syscall.SIGUSR2:
			// reload
			base.Info("Reload")
		default:
			// stop
			base.Info("Stop")
			handler.Stop()
			return
		}
	}
}
