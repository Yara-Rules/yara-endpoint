package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Yara-Rules/yara-endpoint/common"
	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
)

/*
1- Comprobar que tiene el fichero/clave de registro.
    1.1- Crear fichero/clave de registro.
2- Enviar _keep alive_
    2.1- Comprobar si hay retorno de comando
    2.2- Ejecutar comando
3- Volver al punto 2
*/

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

	// showVersion show up the version
	ShowVersion = false
)

const NUM_REGISTRY_TRIES = 3

func printVerion() {
	fmt.Fprintf(os.Stdout, "Yara-Endpoint %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build ID %s\n", BuildID)
	fmt.Fprintf(os.Stdout, "Build on %s\n", BuildDate)
	os.Exit(0)
}

func init() {
	flag.StringVar(&Server, "server", Server, "Server IP/DNS")
	flag.IntVar(&Port, "port", Port, "Server port")
	flag.BoolVar(&ShowVersion, "version", ShowVersion, "Show version")

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	flag.Parse()

	validateFlags()

	log.Info("*** Starting Yara-Endpint ***")

	nc := NewClient(Server, strconv.Itoa(Port))
	log.Debug("Created new Client")

	if err := nc.CheckRegister(); err != nil {
		log.Fatal("Incorrect config file format")
	} else {
		if nc.ULID == "" {
			log.Info("Endpoint no registered. Registering...")

			msg := new(common.Message)
			carryOn := false

			for i := 1; i <= NUM_REGISTRY_TRIES; i++ {
				log.Infof("Sending register command %d/%d.", i, NUM_REGISTRY_TRIES)
				msg, err = nc.Register()
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

				nc.SaveConfig(msg.ULID)
			} else {
				log.Fatal("Reached max tries for the registry process.")
			}
		}

		// Endpoint registered and configured properly

		for {
			log.Info("Sending ping command")
			if msg, err := nc.SendPing(); err != nil {
				log.Error(err)
			} else {
				log.Info("Get pong")
				pp.Println(msg)
			}
			time.Sleep(5 * time.Second)
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
