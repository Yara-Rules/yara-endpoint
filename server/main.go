package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Yara-Rules/yara-endpoint/server/config"
	"github.com/Yara-Rules/yara-endpoint/server/database"
	"github.com/Yara-Rules/yara-endpoint/server/router"
	log "github.com/sirupsen/logrus"
)

var (
	// Version is the bot version
	Version = "v0.0.0"

	// BuildID is the build identifier
	BuildID = ""

	// BuildDate is the compilation date
	BuildDate = ""

	// Log stuff
	// logLevel
	logLevel = "info"
	// logOutput
	logOutput = "-"
	// logFormat
	logFormat = "json"

	configFileName = "yara-server.ini"

	// ShowVersion show up the version
	ShowVersion = false

	DB *database.DataStore
)

func printVerion() {
	fmt.Fprintf(os.Stdout, "Yara-Endpoint %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build ID %s\n", BuildID)
	fmt.Fprintf(os.Stdout, "Build on %s\n", BuildDate)
	os.Exit(0)
}

func init() {
	flag.StringVar(&configFileName, "configFile", configFileName, "Configuration file")
	flag.StringVar(&logLevel, "logLevel", logLevel, "Log level")
	flag.StringVar(&logOutput, "logOutput", logOutput, "Log output file")
	flag.StringVar(&logFormat, "logFormat", logFormat, "Log output format")
	flag.BoolVar(&ShowVersion, "version", ShowVersion, "Show version")
}

func main() {
	// Parse command line flags
	flag.Parse()
	setLog()
	log.Infof("** Yara-Endpoint Server %s **", Version)

	// Read config
	config.LoadConfig(configFileName)

	// Setting up TCP Server
	log.Info("Starting TCP Server")
	srv := NewServer()
	srv.Connect(config.CFG)
	defer srv.Close()
	log.Infof("TCP Server started and listening on %s:%d", config.CFG.TCPServer.Addr, config.CFG.TCPServer.Port)

	// Setting up WEB Server
	log.Info("Starting WEB Server")
	m := router.NewMacaron()
	router.RegisterRoutes(m)

	// Forking WEB Server
	go m.Run(config.CFG.Webserver.Host, config.CFG.Webserver.Port)
	log.Infof("WEB Server started and listening on %s:%d", config.CFG.Webserver.Host, config.CFG.Webserver.Port)

	DB = database.NewDataStore(config.CFG)
	defer DB.Close()

	// Start TCP main loop
	log.Info("Waiting for connections...")
	for {

		conn, err := srv.Accept()
		if err != nil {
			continue
		}
		log.Debugf("New connection accepted from %s", conn.RemoteAddr())
		go srv.HandleClient()
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
