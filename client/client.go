package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
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
	log.Debug("Connecting to the server...")
	conn, err := net.Dial("tcp", c.connString)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Connected to the server")
	c.conn = conn
	c.w = gob.NewEncoder(conn)
	c.r = gob.NewDecoder(conn)
}

// close terminates a connection with Yara-Endpoint server.
func (c *Client) close() {
	c.Lock()
	defer c.Unlock()
	c.conn.Close()
	log.Debug("Closed connection to the server")
}

// send sends the message to the Yara-Endpoint server.
func (c *Client) send(msg *message.Message) error {
	c.Lock()
	defer c.Unlock()

	log.Debug("Sending message to the server")
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

	log.Debug("Receiving message from the server")
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

	log.Debug("Receiving and populating message from the server")
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

	log.Debug("Sending and receiving from the server")
	c.send(msg)
	return c.receiveWithMessage(msg)
}

// IsRegistered checks whether the Endpoint is registered or not.
// It returns the configuration either fill it up or empty.
func (c *Client) IsRegistered() error {
	log.Info("Checking whether endpoint is registered")
	v := new(Config)

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Warn("Configuration file not found.")
		return nil
	}

	if err := ini.MapTo(v, configFileName); err != nil {
		log.Warn("Configuration file does not match with configuration structure.")
		return err
	}

	log.Infof("Endpoint registered with ULID: %s", v.ULID)
	c.ULID = v.ULID
	return nil
}

// register performs the registation and it returns the message received from the server.
func (c *Client) register() (*message.Message, error) {
	log.Debug("Registering endpoint")
	m := message.NewMessage()

	hostname, err := os.Hostname()
	if err != nil {
		log.Warn("Unable to get the endpoint hostname")
		return m, err
	}

	m.CMD = command.Register
	m.Data = fmt.Sprintf("%s|%s", hostname, Version)

	if err := c.sendAndReceive(m); err != nil {
		return message.NewMessage(), err
	}

	log.Infof("Endpoint registered with ULID: %s", m.ULID)
	c.ULID = m.ULID

	return m, nil
}

// saveConfig writes the configuration to a file
func (c *Client) saveConfig(ulid string) error {
	log.Debug("Saving configuration")
	v := new(Config)

	v.ULID = ulid

	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, v)
	if err != nil {
		log.Warn("Unable to generate a valid configuration")
		return err
	}

	log.Infof("Saving configuration to %s", configFileName)
	if runtime.GOOS == "windows" {
		file, err := os.Create(configFileName)
		if err != nil {
			log.Errorf("Unable to create the configuration file. Report your ULID: %s", ulid)
			os.Exit(-1)
		}
		file.Close()
	}
	cfg.SaveTo(configFileName)

	return nil
}

// sendPing
func (c *Client) sendPing() (*message.Message, error) {
	log.Debug("Sending ping to the server")
	m := message.NewMessage()

	m.ULID = c.ULID
	m.CMD = command.Ping
	m.ResultYara = make(map[string][]yara.MatchRule)

	err := c.sendAndReceive(m)

	return m, err
}

func (c *Client) validateMsgFields(msg *message.Message) {
	log.Debug("Validating invalid fields in the message")
	if msg.TaskID == "" {
		log.Warn("Field TaskID is empty, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.TaskIDNotProvided]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
	if msg.Params == "" {
		log.Warn("Field Param is empty, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.ParamsNotProvided]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
	if msg.Data == "" {
		log.Warn("Field Data is empty, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.DataNotProvided]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
	}
}

func (c *Client) checkCompilerErr(msg *message.Message, err error) bool {
	if err != nil {
		log.Warnf("Yara compiler says %s, sending error back", err)
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToGetYaraCompiler]
		msg.ErrorID = errors.UnableToGetYaraCompiler
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (c *Client) checkGetRulesErr(msg *message.Message, err error) bool {
	if err != nil {
		log.Warnf("Yara cannot get rules because %s, sending error back", err)
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToGetRulesFromYaraCompiler]
		msg.ErrorID = errors.UnableToGetRulesFromYaraCompiler
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (c *Client) checkFileExistErr(msg *message.Message, err error) bool {
	if os.IsNotExist(err) {
		log.Warn("Unable to find the file, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.FileDoesNotExist]
		msg.ErrorID = errors.FileDoesNotExist
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (c *Client) checkScanErr(msg *message.Message, err error, errorID errors.Error) bool {
	if err != nil {
		log.Warn("Error while scanning, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errorID]
		msg.ErrorID = errorID
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (c *Client) checkMsgErr(msg *message.Message, err error) bool {
	if err != nil && fmt.Sprintf("%s", err) != "EOF" {
		log.Warn("Got an empty or EOF message from server, sending error back")
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.SendingMsg]
		msg.ErrorID = errors.SendingMsg
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return true
	}
	return false
}

func (c *Client) sendAndPrintMsg(msg *message.Message) {
	log.Debug("Sending message to the server and print it out")
	err := c.sendAndReceive(msg)
	if err != nil {
		log.Warnf("Unable to send the message to the server, err: %s", err)
	}
}

// ScanFile
func (c *Client) ScanFile(msg *message.Message) {
	log.Debug("Running ScanFile routine.")

	c.validateMsgFields(msg)

	comp, err := yara.NewCompiler()
	if c.checkCompilerErr(msg, err) {
		return
	}

	comp.AddString(msg.Data, "")
	rules, err := comp.GetRules()
	if c.checkGetRulesErr(msg, err) {
		return
	}

	_, err = os.Stat(msg.Params)
	if c.checkFileExistErr(msg, err) {
		return
	}

	matches, err := rules.ScanFile(msg.Params, 0, SCAN_TIMEOUT*time.Second)
	if c.checkScanErr(msg, err, errors.ScanningFile) {
		return
	}
	log.Debug("ScanFile finished")

	msg.ResultYara[msg.Params] = matches

	log.Debug("Sending results for ScanFile to server")

	err = c.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if c.checkMsgErr(msg, err) {
		return
	}
}

// ScanDir
func (c *Client) ScanDir(msg *message.Message) {
	log.Debug("Running ScanDir routine.")

	c.validateMsgFields(msg)

	comp, err := yara.NewCompiler()
	if c.checkCompilerErr(msg, err) {
		return
	}

	comp.AddString(msg.Data, "")
	rules, err := comp.GetRules()
	if c.checkGetRulesErr(msg, err) {
		return
	}

	_, err = os.Stat(msg.Params)
	if c.checkFileExistErr(msg, err) {
		return
	}

	fileList := make([]string, 0)
	err = filepath.Walk(msg.Params, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return err
	})
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.UnableToReadFilesInFolder]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	for _, file := range fileList {
		fileMatch, err := rules.ScanFile(file, 0, SCAN_TIMEOUT*time.Second)
		if c.checkScanErr(msg, err, errors.ScanningDir) {
			return
		}
		msg.ResultYara[file] = fileMatch
	}

	log.Debug("ScanDir finished")
	log.Debug("Sending results for ScanDir to server")

	err = c.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if c.checkMsgErr(msg, err) {
		return
	}
}

// ScanPID
func (c *Client) ScanPID(msg *message.Message) {
	log.Debug("Running ScanPID routine.")

	c.validateMsgFields(msg)

	comp, err := yara.NewCompiler()
	if c.checkCompilerErr(msg, err) {
		return
	}

	comp.AddString(msg.Data, "")
	rules, err := comp.GetRules()
	if c.checkGetRulesErr(msg, err) {
		return
	}

	pid, err := strconv.Atoi(msg.Params)
	if err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.BadParams]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	if _, err := os.FindProcess(pid); err != nil {
		// Error
		msg.Error = true
		msg.ErrorMsg = errors.Errors[errors.PIDProcessNotFound]
		c.sendAndPrintMsg(msg)
		// TODO: Check the response, it may contain what to do next
		return
	}

	matches, err := rules.ScanProc(pid, 0, SCAN_TIMEOUT*time.Second)
	fmt.Printf("ERR: %v", err)
	if c.checkScanErr(msg, err, errors.ScanningPID) {
		return
	}
	log.Debug("ScanPID finished")

	msg.ResultYara[msg.Params] = matches

	log.Debug("Sending results for ScanPID to server")

	err = c.sendAndReceive(msg)

	// TODO: Process de response
	// This should be reviewed, because sending message after getting an error by sending a message
	// it couldn't be the best idea
	if c.checkMsgErr(msg, err) {
		return
	}
}

// Ping
func (c *Client) Ping() *message.Message {
	log.Info("Sending PING command")

	msg, err := c.sendPing()
	if err != nil {
		log.Error(err)
	}
	return msg
}

// RegisterEndpoint
func (c *Client) RegisterEndpoint() {
	log.Debug("Running RegisterEndpoint routine.")
	var err error
	msg := new(message.Message)
	carryOn := false

	for i := 1; i <= NUM_REGISTRY_TRIES; i++ {
		log.Infof("Sending <%s> command %d/%d.", command.Alias[command.Register], i, NUM_REGISTRY_TRIES)
		msg, err = c.register()
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

		c.saveConfig(msg.ULID)
	} else {
		log.Debug("Reached max tries for the registry process.")
		log.Fatal("Reached max tries for the registry process.")
	}
}
