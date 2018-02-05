package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/Yara-Rules/yara-endpoint/server/config"
)

type Server struct {
	conn net.Conn
	w    *gob.Encoder
	r    *gob.Decoder
	db   *mgo.Database
	L    net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Connect(c *config.CFG) {
	// Setting up TCP Server
	log.Info("Starting TCP Server")
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.CFG.TCPServer.Addr, c.CFG.TCPServer.Port))
	checkErr(err)

	s.L, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	log.Infof("TCP Server started and listening on %s:%d", c.CFG.TCPServer.Addr, c.CFG.TCPServer.Port)
}

func (s *Server) Close() {
    s.L.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
