package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var host = flag.String("host", "127.0.0.1", "host")
var port = flag.String("port", "6565", "port")
var count = flag.Int("count", 10, "count")

func main() {
	flag.Parse()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGQUIT)
	go func() {
		for i := 0; i < *count; i++ {
			go func() {
				conn, err := net.DialTimeout("tcp", *host+":"+*port, time.Second*60)
				if err != nil {
					logrus.Info("Error connecting:", err)
					os.Exit(1)
				}
				logrus.Info("Connecting to " + *host + ":" + *port)
				handle(conn)
			}()
			//time.Sleep(time.Second * 2)
		}
	}()
	<-c
	logrus.Info("exit .....")
}
func handle(conn net.Conn) {
	defer conn.Close()
	count := 0
	for {
		time.Sleep(time.Second * 5)
		_, err := conn.Write([]byte("hello serverk"))
		if err != nil {
			logrus.Error("Error to send message because of ", err.Error())
			logrus.Info("关闭连接.....")
			break
		}
		count++
		if count > 5 {
			logrus.Info("关闭连接.....")
			break
		}
	}
}
