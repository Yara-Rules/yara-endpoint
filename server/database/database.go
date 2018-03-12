package database

import (
	"fmt"

	"github.com/Yara-Rules/yara-endpoint/server/config"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type DataStore struct {
	host    string
	port    int
	dbName  string
	session *mgo.Session
}

func NewDataStore(c *config.Config) *DataStore {
	sess, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%d", c.Database.Server, c.Database.Port))
	if err != nil {
		log.Error("Unable to connect to the database.")
		log.Fatalf("Unable to connect to the database. Database says: %v", err)
	}
	sess.SetSafe(&mgo.Safe{})
	sess.SetMode(mgo.Monotonic, true)
	return &DataStore{
		host:    c.Database.Server,
		port:    c.Database.Port,
		dbName:  c.Database.DBName,
		session: sess,
	}
}

func (d *DataStore) NewDataStore() *DataStore {
	return &DataStore{
		host:    d.host,
		port:    d.port,
		dbName:  d.dbName,
		session: d.session.Copy(),
	}
}

func (d *DataStore) C(c string) *mgo.Collection {
	return d.session.DB(d.dbName).C(c)
}

func (d *DataStore) Close() {
	d.session.Close()
}

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
