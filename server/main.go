package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"
	"github.com/Yara-Rules/yara-endpoint/common/message"
	"github.com/Yara-Rules/yara-endpoint/server/config"
	"github.com/Yara-Rules/yara-endpoint/server/database"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	"github.com/Yara-Rules/yara-endpoint/server/router"
	"github.com/k0kubun/pp"
	"github.com/oklog/ulid"
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

// var numMessages = 1
var i = 0

var (
	cfile string = "yesconf.ini"
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.CFG.TCPServer.Addr, config.CFG.TCPServer.Port))
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	defer listener.Close()
	log.Infof("TCP Server started and listening on %s:%d", config.CFG.TCPServer.Addr, config.CFG.TCPServer.Port)

	// Setting up WEB Server
	log.Info("Starting WEB Server")
	m := router.NewMacaron()
	router.RegisterRoutes(m)

	// Setting up Mongo Connection
	sess, err := database.GetSession()
	checkErr(err)
	sess.SetSafe(&mgo.Safe{})
	db := database.GetDb(sess)

	// Forking WEB Server
	go m.Run(config.CFG.Webserver.Host, config.CFG.Webserver.Port)
	log.Infof("WEB Server started and listening on %s:%d", config.CFG.Webserver.Host, config.CFG.Webserver.Port)

	// Start TCP main loop
	log.Info("Waiting for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Debugf("New connection accepted from %s", conn.RemoteAddr())
		go handleClient(conn, db)
	}
}

func handleClient(conn net.Conn, db *mgo.Database) {
	defer conn.Close()

	var msg message.Message
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	dec.Decode(&msg)

	if msg.CMD == command.Reserved {
		return
	}

	// fmt.Printf("GET: ")
	// fmt.Printf("(MSG: %d) Received message from %s. Which issed a command: \"%s\"\n", numMessages, msg.ULID, command.Alias[msg.CMD])
	// pp.Println(msg)
	// fmt.Println("---------------------------------------")
	// numMessages++

	switch msg.CMD {
	case command.Register:
		log.Info("Processing Register command")
		// New endpoint
		id := newULID()

		// TODO: Check Hostname value

		ep := &models.Endpoint{
			ULID:     id.String(),
			Hostname: msg.Data,
			Tags:     []string{},
			LastPing: time.Now(),
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}

		db.C(models.Endpoints).Insert(ep)
		log.Debug("Endpoint inserted")

		msg.ULID = id.String()

		enc.Encode(msg)

	case command.Ping:
		log.Info("Processing Ping command")
		err := db.C(models.Endpoints).Update(bson.M{"ulid": msg.ULID}, bson.M{"$set": bson.M{"last_ping": time.Now()}})
		if err == mgo.ErrNotFound {
			log.Errorf("Sending PONG with errors")
			msg.Result = "Pong"
			msg.Error = true
			msg.ErrorID = errors.NeedsRegister
			msg.ErrorMsg = errors.Errors[errors.NeedsRegister]
			enc.Encode(msg)
			return
		} else {
			if jobs, pending := pendingJobs(db, msg); pending {
				now := time.Now().Unix()
				for idx, job := range jobs.Tasks {
					if now >= job.When.Unix() {
						msg.TaskID = job.TaskID
						msg.CMD = job.Command
						msg.Params = job.Params
						msg.Data = extractRules(db, job.Rules)
						remainingJobs := append(jobs.Tasks[:idx], jobs.Tasks[idx+1:]...)
						if len(remainingJobs) == 0 {
							err := db.C(models.Schedules).Remove(bson.M{"_id": jobs.ID})
							if err == mgo.ErrNotFound {
								log.Errorf("Unable to remove Schedule for ID %s", jobs.ID)
								// Due to the error we will wait for the next ping
								msg.Result = "Pong"
								msg.TaskID = ""
								msg.Params = ""
								msg.Data = ""
								msg.Error = true
								msg.ErrorID = errors.UnableToUpdateDB
								msg.ErrorMsg = errors.Errors[errors.UnableToUpdateDB]
							}
						} else {
							err := db.C(models.Schedules).Update(bson.M{"_id": jobs.ID}, bson.M{"$set": bson.M{"tasks": remainingJobs}})
							if err == mgo.ErrNotFound {
								log.Errorf("Unable to update Schedule for ID %s", jobs.ID)
								// Due to the error we will wait for the next ping
								msg.Result = "Pong"
								msg.TaskID = ""
								msg.Params = ""
								msg.Data = ""
								msg.Error = true
								msg.ErrorID = errors.UnableToUpdateDB
								msg.ErrorMsg = errors.Errors[errors.UnableToUpdateDB]
							}
						}
						pp.Println(msg)
						enc.Encode(msg)
						// Finishing the loop
						break
					}
				}
			} else {
				msg.Result = "Pong"
				enc.Encode(msg)
			}
		}

	case command.ScanFile:
		// TODO:
		// - check error
		// - save the result
		// - send a message back
		fmt.Println("ScanFile result")
		pp.Println(msg)
		msg.Result = "OK"
		enc.Encode(msg)

	case command.ScanDir:
		// TODO:
		// - check error
		// - save the result
		// - send a message back
		fmt.Println("ScanDir result")
		pp.Println(msg)
		msg.Result = "OK"
		enc.Encode(msg)

	case command.ScanPID:
		// TODO:
		// - check error
		// - save the result
		// - send a message back
		fmt.Println("ScanPID result")
		pp.Println(msg)
		msg.Result = "OK"
		enc.Encode(msg)
	}
}

func newULID() ulid.ULID {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func pendingJobs(db *mgo.Database, msg message.Message) (*models.Schedule, bool) {
	schedule := new(models.Schedule)
	err := db.C(models.Schedules).Find(bson.M{"ulid": msg.ULID}).One(schedule)
	if err == mgo.ErrNotFound {
		log.Debugf("No pending jobs for %s", msg.ULID)
		return nil, false
	}
	log.Debugf("Pending jobs for %s", msg.ULID)
	return schedule, true
}

func extractRules(db *mgo.Database, RuleList []bson.ObjectId) string {
	var rules string = ""
	var rule models.Rule
	for _, ruleID := range RuleList {
		err := db.C(models.Rules).Find(bson.M{"_id": ruleID}).One(&rule)
		if err != nil {
			log.Errorf("Rule %s not found", ruleID)
			continue
		}
		rules += rule.Data
	}
	return rules
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
