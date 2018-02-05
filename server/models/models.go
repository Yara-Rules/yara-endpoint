package models

import (
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"

	"gopkg.in/mgo.v2/bson"
)

const (
	Endpoints = "endpoints"
	Schedules = "schedules"
	Rules     = "rules"
)

type State int

const (
	Initial State = iota
	Running
	Finished
)

var States = map[State]string{
	Initial:  "Initial",
	Running:  "Running",
	Finished: "Finished",
}

type Endpoint struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ULID     string        `bson:"ulid"`
	Hostname string        `bson:"hostame"`
	Tags     []string      `bson:"tags"`
	LastPing time.Time     `bson:"last_ping"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

type Schedule struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ULID     string        `bson:"ulid"`
	Tasks    []Task        `bson:"tasks"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}

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

type Rule struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Tags     []string      `bson:"tags"`
	Data     string        `bson:"data"`
	CreateAt time.Time     `bson:"created_at"`
	UpdateAt time.Time     `bson:"updated_at"`
}
