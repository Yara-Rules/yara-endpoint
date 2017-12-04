package main

type Command int

const (
	Register Command = iota
	Ping
	Scan
	ScanFile
	ScanMemory
	ScanPID
)

type Message struct {
	ULID     string
	CMD      Command
	Params   string
	Result   string
	Data     string
	Error    bool
	ErrorMsg string
}

func NewMessage() *Message {
	return &Message{}
}
