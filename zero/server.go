package zero

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	addr              string
	port              int16
	associatedManager *AssociatedManager
	timeWheel         *TimeWheel
	listener          *net.TCPListener
}

var timeWheel *TimeWheel = nil

func NewServer(addr string, port int16) *Server {
	t := NewTimeWheel(time.Second*1, 60, func(e SlotElement) {
		c := e.(*Conn)
		c.Close()
	})
	timeWheel = t
	return &Server{
		addr:              addr,
		port:              port,
		associatedManager: NewAssociatedManager(),
		timeWheel:         t,
	}
}

func (s *Server) ListenAndServer() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.addr, s.port))
	if err != nil {
		logrus.Error(err)
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Info("start accept tcp .....")
	s.timeWheel.Start()

	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			logrus.Error(err)
			continue
		}

		c := NewConn(tcpConn, s.timeWheel, func(c *Conn) {
			s.associatedManager.Del(c.GetID())
		})
		s.associatedManager.Set(c.GetID(), c)
		logrus.Info("accept >>>  ", c.GetID())
		go c.Start()
	}
	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}
