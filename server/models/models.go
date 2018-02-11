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
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ULID     string        `bson:"ulid"`
	Hostname string        `bson:"hostame"`
	Tags     []string      `bson:"tags"`
	LastPing time.Time     `bson:"last_ping"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

/* Schematic for schedule collection */

// Schedule collection
type Schedule struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ULID     string        `bson:"ulid"`
	Tasks    []Task        `bson:"tasks"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

//Task type
type Task struct {
	TaskID   string          `bson:"task_id"`
	Command  command.Command `bson:"command"`
	Rules    []bson.ObjectId `bson:"rules"`
	Params   string          `bson:"params"`
	When     time.Time       `bson:"when"`
	Status   State           `bson:"status"`
	CreateAt time.Time       `bson:"created_at"`
	UpdateAt time.Time       `bson:"updated_at"`
}

/* Schematic for rule collection */

// Rule collection
type Rule struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Tags     []string      `bson:"tags"`
	Data     string        `bson:"data"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

/* Schemtic for reports collection */

// Report collection
type Report struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ULID     string        `bson:"ulid"`
	Reports  []LReports    `bson:"reports"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

type LReports struct {
	ReportID string    `bson:"report_id"`
	Task     Task      `bson:"task"`
	Result   []Result  `bson:"result"`
	CreateAt time.Time `bson:"created_at"`
	UpdateAt time.Time `bson:"updated_at"`
}

type Result struct {
	File      string                 `bson:"file"`
	RuleName  string                 `bson:"rule_name"`
	Namespace string                 `bson:"namespace"`
	Tags      []string               `bson:"tags"`
	Meta      map[string]interface{} `bson:"meta"`
	Strings   []String               `bson:"strings"`
}

type String struct {
	Name   string `bson:"name"`
	Offset uint64 `bson:"offset"`
}

/* Schematic for error collection */

// Error collection
type Error struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	ULID        string        `bson:"ulid"`
	TaskID      string        `bson:"task_id"`
	ErrorID     errors.Error  `bson:"error_id"`
	ErrorMsg    string        `bson:"error_msg"`
	Acknowledge bool          `bson:"acknowledge"`
	CreateAt    time.Time     `bson:"created_at"`
}
