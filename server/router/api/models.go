package api

import (
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	"gopkg.in/mgo.v2/bson"
)

type error_ struct {
	Error    bool   `json:"error"`
	ErrorID  int    `json:"error_id"`
	ErrorMsg string `json:"error_msg"`
}

type dashboard struct {
	Asset []models.Endpoint `json:"assets"`
	Rules []models.Rule     `json:"rules"`
}

type publicTasks struct {
	ULID     string `json:"ulid"`
	Hostname string `json:"hostname"`
	Task     struct {
		Command  command.Command `json:"command"`
		Rules    []string        `json:"rules"`
		When     time.Time       `json:"when"`
		Status   models.State    `json:"status"`
		UpdateAt time.Time       `json:"updated_at"`
	} `json:"task"`
}

// Task states
type state int

// Report collection
type report struct {
	ID       bson.ObjectId `json:"-"`
	ULID     string        `json:"ulid"`
	Hostname string        `json:"hostname"`
	Reports  []struct {
		ReportID string `json:"report_id"`
		Task     struct {
			TaskID   string          `json:"task_id"`
			Command  command.Command `json:"command"`
			Rules    []bson.ObjectId `json:"rules"`
			Params   string          `json:"params"`
			When     time.Time       `json:"when"`
			Status   state           `json:"status"`
			CreateAt time.Time       `json:"created_at"`
			UpdateAt time.Time       `json:"updated_at"`
		} `json:"task"`
		Result []struct {
			File      string                 `json:"file"`
			RuleName  string                 `json:"rule_name"`
			Namespace string                 `json:"namespace"`
			Tags      []string               `json:"tags"`
			Meta      map[string]interface{} `json:"meta"`
			Strings   []struct {
				Name   string `json:"name"`
				Offset uint64 `json:"offset"`
			} `json:"strings"`
		} `json:"result"`
		CreateAt time.Time `json:"created_at"`
		UpdateAt time.Time `json:"updated_at"`
	} `json:"reports"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

// type lreports struct {
// 	ReportID string    `bson:"report_id"   json:"report_id"`
// 	Task     task      `bson:"task"        json:"task"`
// 	Result   []result  `bson:"result"      json:"result"`
// 	CreateAt time.Time `bson:"created_at"  json:"created_at"`
// 	UpdateAt time.Time `bson:"updated_at"  json:"updated_at"`
// }

//Task type
// type task struct {
// 	TaskID   string          `bson:"task_id"     json:"task_id"`
// 	Command  command.Command `bson:"command"     json:"command"`
// 	Rules    []bson.ObjectId `bson:"rules"       json:"rules"`
// 	Params   string          `bson:"params"      json:"params"`
// 	When     time.Time       `bson:"when"        json:"when"`
// 	Status   state           `bson:"status"      json:"status"`
// 	CreateAt time.Time       `bson:"created_at"  json:"created_at"`
// 	UpdateAt time.Time       `bson:"updated_at"  json:"updated_at"`
// }

// type result struct {
// 	File      string                 `bson:"file"       json:"file"`
// 	RuleName  string                 `bson:"rule_name"  json:"rule_name"`
// 	Namespace string                 `bson:"namespace"  json:"namespace"`
// 	Tags      []string               `bson:"tags"       json:"tags"`
// 	Meta      map[string]interface{} `bson:"meta"       json:"meta"`
// 	Strings   []string_              `bson:"strings"    json:"strings"`
// }

// type string_ struct {
// 	Name   string `bson:"name"    json:"name"`
// 	Offset uint64 `bson:"offset"  json:"offset"`
// }
