package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"
	"github.com/Yara-Rules/yara-endpoint/common/message"
	"github.com/go-ini/ini"
	yara "github.com/hillu/go-yara"
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
func NewClient(ip, port string) *Client {
	return &Client{
		connString: ip + ":" + port,
	}
}

// connect establishes a connection to a especific Yara-Endpoint server.
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

// close terminates a connection with Yara-Endpoint server.
func (c *Client) close() {
	c.Lock()
	defer c.Unlock()
	c.conn.Close()
}

// send sends the message to the Yara-Endpoint server.
func (c *Client) send(msg *message.Message) error {
	c.Lock()
	defer c.Unlock()

	if err := c.w.Encode(msg); err != nil {
		return err
	}

	return nil
}

// receive gets a message back from the Yara-Endpoint server.
func (c *Client) receive() (*message.Message, error) {
	c.Lock()
	defer c.Unlock()

	msg := message.NewMessage()

	if err := c.r.Decode(&msg); err != nil {
		return msg, err
	}

	return msg, nil
}

// receiveWithMessage gets a message back from the Yara-Endpoint server.
// Moreover it populates `msg` with the response data.
func (c *Client) receiveWithMessage(msg *message.Message) error {
	c.Lock()
	defer c.Unlock()

	if err := c.r.Decode(&msg); err != nil {
		return err
	}

	return nil
}

// sendAndReceive sends the message to the Yara-Endpoint server.
// Moreover it populates `msg` with the response data.
func (c *Client) sendAndReceive(msg *message.Message) error {
	c.connect()
	defer c.close()

	c.send(msg)
	return c.receiveWithMessage(msg)
}

// IsRegistered checks whether the Endpoint is registered or not.
// It returns the configuration either fill it up or empty.
func (c *Client) IsRegistered() error {
	v := new(Config)

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Warn("Configuration file not found.")
		return nil
	}

	if err := ini.MapTo(v, configFileName); err != nil {
		return err
	}

	c.ULID = v.ULID
	return nil
}

// register performs the registation and it returns the message received from the server.
func (c *Client) register() (*message.Message, error) {
	m := message.NewMessage()

	hostname, err := os.Hostname()
	if err != nil {
		return m, err
	}

	m.CMD = command.Register
	m.Data = hostname

	if err := c.sendAndReceive(m); err != nil {
		return message.NewMessage(), err
	}

	c.ULID = m.ULID

	return m, nil
}

// saveConfig writes the configuration to a file
func (c *Client) saveConfig(ulid string) error {
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

// sendPing
func (c *Client) sendPing() (*message.Message, error) {
	m := message.NewMessage()

	m.ULID = c.ULID
	m.CMD = command.Ping
	m.ResultYara = make(map[string][]yara.MatchRule)

	err := c.sendAndReceive(m)

	return m, err
}

func (nc *Client) validateMsgFields(msg *message.Message) {
	if msg.TaskID == "" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.TaskIDNotProvided]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
	if msg.Params == "" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.ParamsNotProvided]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
	if msg.Data == "" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.DataNotProvided]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
}

func (nc *Client) checkCompilerErr(msg *message.Message, err error) bool {
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToGetYaraCompiler]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (nc *Client) checkGetRulesErr(msg *message.Message, err error) bool {
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToGetRulesFromYaraCompiler]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (nc *Client) sendAndPrintMsg(msg *message.Message) {
	err := nc.sendAndReceive(msg)
	if err != nil {
		log.Println("Error sent and get it back:")
		log.Println("Err: ", err)
	}
}

// ScanFile
func (nc *Client) ScanFile(msg *message.Message) {
	log.Debug("Running ScanFile routine.")

	nc.validateMsgFields(msg)

	c, err := yara.NewCompiler()
	if nc.checkCompilerErr(msg, err) {
		return
	}

	c.AddString(msg.Data, "")
	rules, err := c.GetRules()
	if nc.checkGetRulesErr(msg, err) {
		return
	}

	if _, err := os.Stat(msg.Params); os.IsNotExist(err) {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.FileDoesNotExist]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	matches, err := rules.ScanFile(msg.Params, 0, SCAN_TIMEOUT*time.Second)
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.ScanningFile]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	msg.ResultYara[msg.Params] = matches

	log.Debug("Sending results for ScanFile to server")

	err = nc.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if err != nil && fmt.Sprintf("%s", err) != "EOF" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.SendingMsg]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
}

// ScanDir
func (nc *Client) ScanDir(msg *message.Message) {
	log.Debug("Running ScanDir routine.")

	nc.validateMsgFields(msg)

	c, err := yara.NewCompiler()
	if nc.checkCompilerErr(msg, err) {
		return
	}

	c.AddString(msg.Data, "")
	rules, err := c.GetRules()
	if nc.checkGetRulesErr(msg, err) {
		return
	}

	if _, err := os.Stat(msg.Params); os.IsNotExist(err) {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.FileDoesNotExist]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	fileList := make([]string, 0)
	err = filepath.Walk(msg.Params, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToReadFilesInFolder]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	for _, file := range fileList {
		fileMatch, err := rules.ScanFile(file, 0, SCAN_TIMEOUT*time.Second)
		if err != nil {
			// Error
			msg.Error = true
			msg.ErrorMsg = errors.Errors[errors.ScanningDir]
			nc.sendAndPrintMsg(msg)
			// TODO: Check the response, it may contain what to do next
			return
		}
		msg.ResultYara[file] = fileMatch
	}

	log.Debug("Sending results for ScanDir to server")

	err = nc.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if err != nil && fmt.Sprintf("%s", err) != "EOF" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.SendingMsg]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
}

// ScanPID
func (nc *Client) ScanPID(msg *message.Message) {
	log.Debug("Running ScanPID routine.")

	nc.validateMsgFields(msg)

	c, err := yara.NewCompiler()
	if nc.checkCompilerErr(msg, err) {
		return
	}

	c.AddString(msg.Data, "")
	rules, err := c.GetRules()
	if nc.checkGetRulesErr(msg, err) {
		return
	}

	pid, err := strconv.Atoi(msg.Params)
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.BadParams]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	if _, err := os.FindProcess(pid); err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.PIDProcessNotFound]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	matches, err := rules.ScanProc(pid, 0, SCAN_TIMEOUT*time.Second)
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.ScanningPID]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	msg.ResultYara[msg.Params] = matches

	log.Debug("Sending results for ScanPID to server")

	err = nc.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if err != nil && fmt.Sprintf("%s", err) != "EOF" {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.SendingMsg]
		nc.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
}

// Ping
func (nc *Client) Ping() *message.Message {
	log.Info("Sending PING command")

	msg, err := nc.sendPing()
	if err != nil {
		log.Error(err)
	}
	return msg
}

// RegisterEndpoint
func (nc *Client) RegisterEndpoint() {
	var err error
	msg := new(message.Message)
	carryOn := false

	for i := 1; i <= NUM_REGISTRY_TRIES; i++ {
		log.Infof("Sending REGISTER command %d/%d.", i, NUM_REGISTRY_TRIES)
		msg, err = nc.register()
		if err != nil {
			log.Error(err)
			log.Infof("Unable to get registered: try %d/%d", i, NUM_REGISTRY_TRIES)
		} else {
			carryOn = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if carryOn {
		log.Infof("Endpoint got registered with ID: %s", msg.ULID)
		log.Debugf("Message (recv): %v\n", msg)

		nc.saveConfig(msg.ULID)
	} else {
		log.Debug("Reached max tries for the registry process.")
		log.Fatal("Reached max tries for the registry process.")
	}
}
