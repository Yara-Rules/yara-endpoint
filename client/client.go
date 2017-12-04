package main

import (
	"encoding/gob"
	"net"
	"os"
	"sync"

	"github.com/go-ini/ini"
)

const configFileName = "config.ini"

type Config struct {
	ULID string `ini:"ulid"`
}

// Client is a basic client to the Yara-Endpoint server.
type Client struct {
	conn net.Conn
	w    *gob.Encoder
	r    *gob.Decoder
	sync.Mutex
}

// NewClient returns a Yara-Endpoint client.
func NewClient() *Client {
	return &Client{}
}

// Connect establishes a connection to a Yara-Endpoint server.
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

// Close terminates a connection to Yara-Endpoint.
func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()
	c.conn.Close()
}

// Send sends the message to the server.
func (c *Client) Send(msg *Message) error {
	c.Lock()

	if err := c.w.Encode(msg); err != nil {
		return err
	}

	c.Unlock()

	return nil
}

func (c *Client) CheckRegister() (Config, error) {
	v := new(Config)

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return v, nil
	}

	if err := ini.MapTo(v, configFileName); err != nil {
		return v, err
	}

	return v, nil
}

func (c *Client) Register() Message, error {
	m = new(Message)

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	m.CMD = Register
	m.Data = hostname

	c.Send(m)
	c.r.Decode(&m)

	return m, nil
}

func (c *Client) SaveConfig(ulid string) error {
	v = new(Config)

	v.ULID = ulid

	cfg := ini.Empty()
	err = ini.ReflectFrom(cfg, v)
	if err != nil {
		return err
	}

	cfg.SaveTo(configFileName)

	return nil
}
