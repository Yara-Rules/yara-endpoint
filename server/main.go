package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/Yara-Rules/yara-endpoint/client/config"
	"github.com/oklog/ulid"
)

func main() {

	service := "0.0.0.0:8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var cfg config.Config
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	dec.Decode(&cfg)

	fmt.Printf("GET: %v\n", cfg)

	if cfg.ID == "" {
		// New endpoint
		t := time.Now()
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		id := ulid.MustNew(ulid.Timestamp(t), entropy)
		cfg.ID = id.String()

		enc.Encode(cfg)
		fmt.Printf("OUT: %v\n", cfg)
	} else {
		fmt.Printf("Endpoint already registered with id: %s\n", cfg.ID)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
