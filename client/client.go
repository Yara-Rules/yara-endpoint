package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/Yara-Rules/yara-endpoint/common"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const configFileName = "config.ini"

// Config holds the client's configuration
type Config struct {
	ULID string `ini:"ulid"`
}

// Client is what represents an endpoint.
type Client struct {
	conn net.Conn
	w    *gob.Encoder
	r    *gob.Decoder
	sync.Mutex
	connString string
	Config
}

// NewClient returns a new empty Yara-Endpoint client.
// func NewClient() *Client {
// 	return &Client{}
// }

// NewClient returns a new empty Yara-Endpoint client.
func NewClient(ip, port string) *Client {
	return &Client{
		connString: ip + ":" + port,
	}
}

// Connect establishes a connection to a especific Yara-Endpoint server.
func (c *Client) Connect(netloc string) error {
	conn, err := net.Dial("tcp", netloc)
	if err != nil {
		return err
	}
	c.conn = conn
	c.w = gob.NewEncoder(conn)
	c.r = gob.NewDecoder(conn)

	return nil
}

// Connect establishes a connection to a especific Yara-Endpoint server.
func (c *Client) connect() {
	conn, err := net.Dial("tcp", c.connString)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Connected to server")
	c.conn = conn
	c.w = gob.NewEncoder(conn)
	c.r = gob.NewDecoder(conn)
}

// Close terminates a connection with Yara-Endpoint server.
func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()
	c.conn.Close()
}

// Send sends the message to the Yara-Endpoint server.
func (c *Client) Send(msg *common.Message) error {
	c.Lock()
	defer c.Unlock()

	if err := c.w.Encode(msg); err != nil {
		return err
	}

	return nil
}

// SendAndReceiveWithMessage sends the message to the Yara-Endpoint server.
// Moreover it populates `msg` with the response data.
func (c *Client) SendAndReceiveWithMessage(msg *common.Message) error {
	c.Send(msg)
	return c.ReceiveWithMessage(msg)
}

// Receive gets a message back from the Yara-Endpoint server.
func (c *Client) Receive() (*common.Message, error) {
	c.Lock()
	defer c.Unlock()

	msg := common.NewMessage()

	if err := c.r.Decode(&msg); err != nil {
		return msg, err
	}

	return msg, nil
}

// ReceiveWithMessage gets a message back from the Yara-Endpoint server.
// Moreover it populates `msg` with the response data.
func (c *Client) ReceiveWithMessage(msg *common.Message) error {
	c.Lock()
	defer c.Unlock()

	if err := c.r.Decode(&msg); err != nil {
		return err
	}

	return nil
}

// CheckRegister checks whether the Endpoint is registered or not.
// It returns the configuration either fill it up or empty.
func (c *Client) CheckRegister() error {
	v := new(Config)

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Error("Configuration file not found.")
		return fmt.Errorf("Configuration file not found.")
	}

	if err := ini.MapTo(v, configFileName); err != nil {
		return err
	}

	c.ULID = v.ULID
	return nil
}

// Register performs the registation and it returns the message received from the server.
func (c *Client) Register() (*common.Message, error) {
	m := common.NewMessage()

	hostname, err := os.Hostname()
	if err != nil {
		return m, err
	}

	m.CMD = common.Register
	m.Data = hostname

	c.connect()
	defer c.Close()

	if err := c.SendAndReceiveWithMessage(m); err != nil {
		return common.NewMessage(), err
	}

	return m, nil
}

// SaveConfig writes the configuration to a file
func (c *Client) SaveConfig(ulid string) error {
	v := new(Config)

	v.ULID = ulid

	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, v)
	if err != nil {
		return err
	}

	cfg.SaveTo(configFileName)

	return nil
}

// SendPing
func (c *Client) SendPing() (*common.Message, error) {
	m := common.NewMessage()

	m.ULID = c.ULID
	m.CMD = common.Ping

	c.connect()
	defer c.Close()
	err := c.SendAndReceiveWithMessage(m)

	return m, err
}
