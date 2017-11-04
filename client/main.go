package main

import (
	"flag"
	"fmt"
	"os"
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

func printVerion() {
	fmt.Fprintf(os.Stdout, "Yara-Endpoint %s\n", Version)
	fmt.Fprintf(os.Stdout, "Build ID %s\n", BuildID)
	fmt.Fprintf(os.Stdout, "Build on %s\n", BuildDate)
	os.Exit(0)
}

func init() {
	flag.StringVar(&Server, "server", Server, "Server IP/DNS")
	flag.IntVar(&Port, "port", Port, fmt.Sprintf("Server port (default: %d)", Port))
	flag.BoolVar(&ShowVersion, "version", ShowVersion, "Show version")
}

func main() {
	flag.Parse()

	validateFlags()
}

func validateFlags() {
	if ShowVersion {
		printVerion()
	}
	if Server == "" {
		fmt.Fprintln(os.Stderr, "You must provide a server address.")
	}
	if Port < 1 || Port > 65535 {
		fmt.Fprintln(os.Stderr, "You must provide a valid port number.")
	}
}
