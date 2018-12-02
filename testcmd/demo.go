package main

import (
	"github.com/0x1306a94/go-message-push/zero"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGQUIT)

	start := make(chan struct{})
	server := zero.NewServer("0.0.0.0", 6565)
	defer func() {
		logrus.Info("stop tcp server ...")
		server.Stop()
	}()
	go func() {
		if err := server.ListenAndServer(); err != nil {
			logrus.Error(err)
			start <- struct{}{}
		}
	}()

	go func() {
		select {
		case <-start:
			os.Exit(-1)
		case <-time.After(5):
			logrus.Info("start http debug expvar ....")
		}
		if err := http.ListenAndServe(":6566", nil); err != nil {
			logrus.Error(err)
		}
	}()
	<-c
	logrus.Info("exit ...")
}
