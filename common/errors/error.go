package errors

type Error int

const (
	Reserved Error = iota
	BadTaskID
	TaskIDNotProvided
	BadParams
	ParamsNotProvided
	BadData
	DataNotProvided
	UnableToGetYaraCompiler
	UnableToGetRulesFromYaraCompiler
	ScanningFile
	ScanningDir
	ScanningPID
	SendingMsg
	UnableToReadFilesInFolder
	FileDoesNotExist
	NeedsRegister
	UnableToUpdateDB
	PIDProcessNotFound
)

var Errors = map[Error]string{
	Reserved:                         "Reserved",
	BadTaskID:                        "Bad task ID",
	TaskIDNotProvided:                "Task ID not provided",
	BadParams:                        "Bad Parameters",
	ParamsNotProvided:                "Parameters not provided",
	BadData:                          "Bad data",
	DataNotProvided:                  "Data not provided",
	UnableToGetYaraCompiler:          "Error while getting a Yara compiler",
	UnableToGetRulesFromYaraCompiler: "Error while getting yara rules from compiler",
	ScanningFile:                     "Error while scanning a file",
	ScanningDir:                      "Error while scanning a directory",
	ScanningPID:                      "Error while scanning a PID",
	SendingMsg:                       "Error while sending message to server",
	UnableToReadFilesInFolder:        "Error while reading files in directory",
	FileDoesNotExist:                 "File does not exist",
	NeedsRegister:                    "Register first",
	UnableToUpdateDB:                 "Unable to update DB",
	PIDProcessNotFound:               "PID Process not found",
}
