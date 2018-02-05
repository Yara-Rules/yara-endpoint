package message

import (
	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"
	yara "github.com/hillu/go-yara"
)

type Message struct {
	ULID       string
	TaskID     string
	CMD        command.Command
	Params     string
	Result     string
	ResultYara map[string][]yara.MatchRule
	Data       string
	Error      bool
	ErrorID    errors.Error
	ErrorMsg   string
}

func NewMessage() *Message {
	return &Message{}
}
