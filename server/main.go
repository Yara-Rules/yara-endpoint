package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/Yara-Rules/yara-endpoint/common"
	"github.com/k0kubun/pp"
	"github.com/oklog/ulid"
)

var numMessages = 1

func main() {

	service := "0.0.0.0:8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("New connection")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var msg common.Message
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	dec.Decode(&msg)

	if msg.CMD == common.Reserved {
		return
	}

	fmt.Printf("GET: ")
	fmt.Printf("(MSG: %d) Received message from %s. Which issed a command: %d\n", numMessages, msg.ULID, msg.CMD)
	pp.Println(msg)
	fmt.Println("---------------------------------------")
	numMessages++

	switch msg.CMD {
	case common.Register:
		// New endpoint
		t := time.Now()
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		id := ulid.MustNew(ulid.Timestamp(t), entropy)
		msg.ULID = id.String()

		enc.Encode(msg)

	case common.Ping:
		msg.Result = "Pong"
		enc.Encode(msg)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
