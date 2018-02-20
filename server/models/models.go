package models

import (
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"

	"gopkg.in/mgo.v2/bson"
)

// MongoDB Collections
const (
	Endpoints = "endpoints"
	Schedules = "schedules"
	Rules     = "rules"
	Reports   = "reports"
	Errors    = "errors"
)

// Task states
type State int

const (
	Initial State = iota
	Running
	Finished
	Failed
)

var States = map[State]string{
	Initial:  "Initial",
	Running:  "Running",
	Finished: "Finished",
	Failed:   "Failed",
}

/* Schematic for endpoint collection */

// Endpoint collection
type Endpoint struct {
	ID            bson.ObjectId `bson:"_id,omitempty"  json:"-"`
	ULID          string        `bson:"ulid"           json:"ulid"`
	Hostname      string        `bson:"hostame"        json:"hostname"`
	ClientVersion string        `bson:"client_version" json:"client_version"`
	Tags          []string      `bson:"tags"           json:"tags"`
	LastPing      time.Time     `bson:"last_ping"      json:"last_ping"`
	CreateAt      time.Time     `bson:"created_at"     json:"created_at"`
	UpdateAt      time.Time     `bson:"updated_at"     json:"updated_at"`
}

/* Schematic for schedule collection */

// Schedule collection
type Schedule struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ULID     string        `bson:"ulid"          json:"ulid"`
	Tasks    []Task        `bson:"tasks"         json:"tasks"`
	CreateAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"    json:"updated_at"`
}

//Task type
type Task struct {
	TaskID   string          `bson:"task_id"     json:"task_id"`
	Command  command.Command `bson:"command"     json:"command"`
	Rules    []bson.ObjectId `bson:"rules"       json:"rules"`
	Params   string          `bson:"params"      json:"params"`
	When     time.Time       `bson:"when"        json:"when"`
	Status   State           `bson:"status"      json:"status"`
	CreateAt time.Time       `bson:"created_at"  json:"created_at"`
	UpdateAt time.Time       `bson:"updated_at"  json:"updated_at"`
}

/* Schematic for rule collection */

// Rule collection
type Rule struct {
	ID       bson.ObjectId `bson:"_id,omitempty"  json:"-"`
	Name     string        `bson:"name"           json:"name"`
	Tags     []string      `bson:"tags"           json:"tags"`
	Data     string        `bson:"data"           json:"data"`
	CreateAt time.Time     `bson:"created_at"     json:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"     json:"updated_at"`
}

/* Schemtic for reports collection */

// Report collection
type Report struct {
	ID       bson.ObjectId `bson:"_id,omitempty"  json:"-"`
	ULID     string        `bson:"ulid"           json:"ulid"`
	Reports  []LReports    `bson:"reports"        json:"reports"`
	CreateAt time.Time     `bson:"created_at"     json:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"     json:"updated_at"`
}

type LReports struct {
	ReportID string    `bson:"report_id"   json:"report_id"`
	Task     Task      `bson:"task"        json:"task"`
	Result   []Result  `bson:"result"      json:"result"`
	CreateAt time.Time `bson:"created_at"  json:"created_at"`
	UpdateAt time.Time `bson:"updated_at"  json:"updated_at"`
}

type Result struct {
	File      string                 `bson:"file"       json:"file"`
	RuleName  string                 `bson:"rule_name"  json:"rule_name"`
	Namespace string                 `bson:"namespace"  json:"namespace"`
	Tags      []string               `bson:"tags"       json:"tags"`
	Meta      map[string]interface{} `bson:"meta"       json:"meta"`
	Strings   []String               `bson:"strings"    json:"strings"`
}

type String struct {
	Name   string `bson:"name"    json:"name"`
	Offset uint64 `bson:"offset"  json:"offset"`
}

/* Schematic for error collection */

// Error collection
type Error struct {
	ID          bson.ObjectId `bson:"_id,omitempty"  json:"-"`
	ULID        string        `bson:"ulid"           json:"ulid"`
	TaskID      string        `bson:"task_id"        json:"task_id"`
	ErrorID     errors.Error  `bson:"error_id"       json:"error_id"`
	ErrorMsg    string        `bson:"error_msg"      json:"error_msg"`
	Acknowledge bool          `bson:"acknowledge"    json:"acknowledge"`
	CreateAt    time.Time     `bson:"created_at"     json:"created_at"`
	UpdateAt    time.Time     `bson:"updated_at"     json:"updated_at"`
}
