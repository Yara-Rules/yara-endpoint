package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// Version is the bot version
	Version = "v0.0.0"

	// BuildID is the build identifier
	BuildID = ""

	// BuildDate is the compilation date
	BuildDate = ""

	// Server is the controller
	Server = ""

	// Port is the connection port to the server
	Port = 8080

	// Cinfiguration file
	configFileName = "yara-endpoint.ini"

	// Log stuff
	// logLevel
	logLevel = "info"
	// logOutput
	logOutput = "-"
	// logFormat
	logFormat = "json"

	// ShowVersion show up the version
	ShowVersion = false
)

const (
	WAIT_BETWEEN_PING  = 5
	NUM_REGISTRY_TRIES = 3
	SCAN_TIMEOUT       = 60
)

func printVerion() {
	fmt.Fprintf(os.Stdout, "Yara-Endpoint %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build ID %s\n", BuildID)
	fmt.Fprintf(os.Stdout, "Build on %s\n", BuildDate)
	os.Exit(0)
}

func init() {
	flag.StringVar(&Server, "server", Server, "Server IP/DNS")
	flag.IntVar(&Port, "port", Port, "Server port")
	flag.StringVar(&configFileName, "config", configFileName, "Configuratin file")
	flag.StringVar(&logLevel, "logLevel", logLevel, "Log level")
	flag.StringVar(&logOutput, "logOutput", logOutput, "Log output file")
	flag.StringVar(&logFormat, "logFormat", logFormat, "Log output format")
	flag.BoolVar(&ShowVersion, "version", ShowVersion, "Show version")
}

func main() {
	flag.Parse()
	validateFlags()
	setLog()

	log.Infof("*** Starting Yara-Endpint %s ***", Version)
	nc := NewClient(Server, strconv.Itoa(Port))
	log.Debug("New client created")

	if err := nc.IsRegistered(); err != nil {
		log.Error("Incorrect config file format")
	} else {
		if nc.ULID == "" {
			log.Info("Endpoint no registered. Registering...")
			nc.RegisterEndpoint()
		}

		// Endpoint it now registered and configured properly, carry on.
		for {
			msg := nc.Ping()

			if msg.Error {
				log.Errorf("Server says: %s", errors.Errors[msg.ErrorID])
				switch msg.ErrorID {
				case errors.NeedsRegister:
					nc.RegisterEndpoint()
				case errors.UnableToUpdateDB:

				}
			} else {
				switch msg.CMD {
				case command.ScanFile:
					log.Info("Received ScanFile command.")
					go nc.ScanFile(msg)
				case command.ScanDir:
					log.Info("Received ScanDir command.")
					go nc.ScanDir(msg)
				case command.ScanPID:
					log.Info("Received ScanPID command.")
					go nc.ScanPID(msg)
				case command.Ping:
				}
			}
			time.Sleep(WAIT_BETWEEN_PING * time.Second)
		}
	}
}

func validateFlags() {
	if ShowVersion {
		printVerion()
	}
	if Server == "" {
		fmt.Fprintln(os.Stderr, "You must provide a server address.")
		os.Exit(1)
	}
	if Port < 0 || Port > 65535 {
		fmt.Fprintln(os.Stderr, "You must provide a valid port number.")
		os.Exit(1)
	}
}

func setLog() {
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	if logOutput == "-" {
		log.SetOutput(os.Stdout)
	} else {
		out, err := os.OpenFile(logOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		log.SetOutput(out)
	}
	switch logLevel {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	}
}
