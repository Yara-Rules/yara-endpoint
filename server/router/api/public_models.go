package api

import (
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/server/models"
)

type Error struct {
	Error    bool   `json:"error"`
	ErrorID  int    `json:"error_id"`
	ErrorMsg string `json:"error_msg"`
}

type PublicDashboard struct {
	Asset []models.Endpoint `json:"assets"`
	Rules []models.Rule     `json:"rules"`
}

type PublicTasks struct {
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
type PublicReports struct {
	ULID     string    `json:"ulid"`
	Hostname string    `json:"hostname"`
	Reports  []Report  `json:"reports"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type Report struct {
	ReportID string    `json:"report_id"`
	Task     Task      `json:"task"`
	Result   []Result  `json:"result"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type Task struct {
	TaskID   string          `json:"task_id"`
	Command  command.Command `json:"command"`
	Rules    []string        `json:"rules"`
	Params   string          `json:"params"`
	When     time.Time       `json:"when"`
	Status   models.State    `json:"status"`
	CreateAt time.Time       `json:"created_at"`
	UpdateAt time.Time       `json:"updated_at"`
}

type Result struct {
	File      string                 `json:"file"`
	RuleName  string                 `json:"rule_name"`
	Namespace string                 `json:"namespace"`
	Tags      []string               `json:"tags"`
	Meta      map[string]interface{} `json:"meta"`
	Strings   []YString              `json:"strings"`
}

type YString struct {
	Name   string `json:"name"`
	Offset uint64 `json:"offset"`
}

type NewRuleForm struct {
	Name string   `json:"name" binding:"Required"`
	Tags []string `json:"tags"`
	Data string   `json:"data" binding:"Required"`
}

type NewTaskForm struct {
	Assets  []string `json:"assets"  binding:"Required"`
	Rules   []string `json:"rules"   binding:"Required"`
	Command string   `json:"command" binding:"Required"`
	Target  string   `json:"target"  binding:"Required"`
}
