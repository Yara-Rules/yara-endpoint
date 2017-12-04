package cmd

type Command int

const (
	Register Command = iota
	KeepAlive
)

type CMD struct {
	Cmd  Command
	Args []string
}

func NewCMD(c Command, a []string) *CMD {
	return &CMD{Cmd: c, Args: a}
}
