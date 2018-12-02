package zero

import (
	"bufio"
	"expvar"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

var stats = expvar.NewMap("tcp")
var connectionCount expvar.Int
var start = time.Now()

func calculateUptime() interface{} {
	return time.Since(start).String()
}
func init() {
	stats.Set("connections", &connectionCount)
	expvar.Publish("uptime", expvar.Func(calculateUptime))
}

type Conn struct {
	raw        *net.TCPConn
	timeWheel  *TimeWheel
	disconnect func(c *Conn)
	stopChan   chan struct{}
	stop       bool
}

func NewConn(raw *net.TCPConn, timeWheel *TimeWheel, disconnect func(c *Conn)) *Conn {
	return &Conn{
		raw:        raw,
		timeWheel:  timeWheel,
		disconnect: disconnect,
		stopChan:   make(chan struct{}),
		stop:       false,
	}
}

func (c *Conn) GetID() string {
	return c.raw.RemoteAddr().String()
}

func (c *Conn) Start() {
	id := c.GetID()
	defer func() {
		logrus.Info(id + " : 关闭连接...")
		c.stop = true
		if c.raw != nil {
			c.raw.Close()
		}
		if c.disconnect != nil {
			c.disconnect(c)
		}
		connectionCount.Add(-1)
	}()
	connectionCount.Add(1)
	reader := bufio.NewReader(c.raw)
	c.timeWheel.Add(c)
Exit:
	for {
		select {
		case <-c.stopChan:
			break
		default:
			logrus.Info(id + " : start read ...")
			line, err := reader.ReadString('k')
			if err != nil {
				logrus.Error("read error: ", err)
				break Exit
			} else {
				logrus.Info(id + " : " + line)
				c.timeWheel.Add(c)
			}
		}
	}
	/*
		reader := bufio.NewReader(c.raw)
		for {
			len, err := io.ReadFull(reader, buffer[0:4])
			if err != nil {
				fmt.Println(err)
				continue
			}
			dataLen, err := io.ReadFull(reader, buffer[0:len])
			if err != nil {
				fmt.Println("读取错误")
				return
			}
			fmt.Println("数据长度: ", dataLen)
		}
	*/
}

func (c *Conn) Close() {
	c.stop = true
	c.timeWheel.Remove(c)
	go func(c *Conn) {
		c.stopChan <- struct{}{}
	}(c)
}

func (c *Conn) Closed() bool {
	return c.stop
}
