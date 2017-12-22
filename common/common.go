package common

type Command int

const (
	Reserved Command = iota
	Register
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
