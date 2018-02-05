package config

import (
	"log"

	ini "gopkg.in/ini.v1"
)

var CFG *Config

type Config struct {
	Database  DB        `ini:"Database"`
	Webserver WebServer `ini:"WebServer"`
	TCPServer TCPServer `ini:"TCPServer"`
}

type DB struct {
	Server string `ini:"Server"`
	Port   int    `ini:"Port"`
	DBName string `ini:"DBName"`
}

type WebServer struct {
	Host string `ini:"Hostname"`
	Port int    `ini:"Port"`
}

type TCPServer struct {
	Addr string `ini:"Address"`
	Port int    `ini:"Port"`
}

func LoadConfig(c string) {
	CFG = new(Config)
	err := ini.MapTo(CFG, c)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
