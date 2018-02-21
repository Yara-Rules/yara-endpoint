package api

import (
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/server/models"
)

type error_ struct {
	Error    bool   `json:"error"`
	ErrorID  int    `json:"error_id"`
	ErrorMsg string `json:"error_msg"`
}

type publicDashboard struct {
	Asset []models.Endpoint `json:"assets"`
	Rules []models.Rule     `json:"rules"`
}

type publicTasks struct {
	ULID     string `json:"ulid"`
	Hostname string `json:"hostname"`
	Task     struct {
		TaskID   string          `json:"task_id"`
		Command  command.Command `json:"command"`
		Rules    []string        `json:"rules"`
		When     time.Time       `json:"when"`
		Status   models.State    `json:"status"`
		CreateAt time.Time       `json:"created_at"`
		UpdateAt time.Time       `json:"updated_at"`
	} `json:"task"`
}

// Report collection
type publicReports struct {
	ULID     string     `json:"ulid"`
	Hostname string     `json:"hostname"`
	Reports  []reports_ `json:"reports"`
	CreateAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `json:"updated_at"`
}

type reports_ struct {
	ReportID string     `json:"report_id"`
	Task     task_      `json:"task"`
	Result   []results_ `json:"result"`
	CreateAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `json:"updated_at"`
}

type task_ struct {
	TaskID   string          `json:"task_id"`
	Command  command.Command `json:"command"`
	Rules    []string        `json:"rules"`
	Params   string          `json:"params"`
	When     time.Time       `json:"when"`
	Status   models.State    `json:"status"`
	CreateAt time.Time       `json:"created_at"`
	UpdateAt time.Time       `json:"updated_at"`
}

type results_ struct {
	File      string                 `json:"file"`
	RuleName  string                 `json:"rule_name"`
	Namespace string                 `json:"namespace"`
	Tags      []string               `json:"tags"`
	Meta      map[string]interface{} `json:"meta"`
	Strings   []strings_             `json:"strings"`
}

type strings_ struct {
	Name   string `json:"name"`
	Offset uint64 `json:"offset"`
}
