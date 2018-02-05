package database

import (
	"fmt"

	"github.com/Yara-Rules/yara-endpoint/server/config"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

func GetSession() (*mgo.Session, error) {
	connstr := fmt.Sprintf("mongodb://%s:%d", config.CFG.Database.Server, config.CFG.Database.Port)
	return mgo.Dial(connstr)
}

func GetDb(sess *mgo.Session) *mgo.Database {
	return sess.DB(config.CFG.Database.DBName)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
