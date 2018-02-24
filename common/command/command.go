package command

type Command int

const (
	Reserved Command = iota
	Register
	Ping
	Scan
	ScanFile
	ScanDir
	ScanPID
)

var Alias = map[Command]string{
	Reserved: "Reserved",
	Register: "Register",
	Ping:     "Ping",
	Scan:     "Scan",
	ScanFile: "ScanFile",
	ScanDir:  "ScanDir",
	ScanPID:  "ScanPID",
}

var RevAlias = map[string]Command{
	"Reserved": Reserved,
	"Register": Register,
	"Ping":     Ping,
	"Scan":     Scan,
	"ScanFile": ScanFile,
	"ScanDir":  ScanDir,
	"ScanPID":  ScanPID,
}
