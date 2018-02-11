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

	// ShowVersion show up the version
	ShowVersion = false
)

var (
	cfile string = "yesconf.ini"
	DB    *database.DataStore
)

func printVerion() {
	fmt.Fprintf(os.Stdout, "Yara-Endpoint %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build ID %s\n", BuildID)
	fmt.Fprintf(os.Stdout, "Build on %s\n", BuildDate)
	os.Exit(0)
}

func init() {
	flag.StringVar(&cfile, "configFile", cfile, "Configuration file")
	flag.BoolVar(&ShowVersion, "version", ShowVersion, "Show version")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.InfoLevel) // DebugLevel
	log.SetLevel(log.DebugLevel) // InfoLevel
}

func main() {
	// Parse command line flags
	flag.Parse()

	log.Infof("** Yara-Endpoint Server %s **", Version)

	// Read config
	config.LoadConfig(cfile)

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
