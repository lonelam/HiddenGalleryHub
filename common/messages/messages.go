package messages

import (
	"encoding/json"
)

type MachineInfoMessage struct {
	MachineName string
	PasswdSum   string
}
type RequestFileMessage struct {
	RelativePath string
}

type RequestDirectoryMessage struct {
	RelativePath string
}

type FileEntry struct {
	Name               string
	RelativePath       string
	ParentRelativePath string
}

type DirectoryEntry struct {
	Name               string
	RelativePath       string
	ParentRelativePath string
}

type DirectoryStructureMessage struct {
	FileEntries      []FileEntry
	DirectoryEntries []DirectoryEntry
}

const (
	MessageTypeInit = iota + 100
	MessageTypeMachineInfo
	MessageTypeRequestDirectory
	MessageTypeRequestFile
	MessageTypeFileStart
	MessageTypeFileContent
	MessageTypeFileFinish
	MessageTypeFileException
	MessageTypeDirectoryStructure
)

func ReadMachineInfoFromMessage(data []byte) (*MachineInfoMessage, error) {
	var v MachineInfoMessage
	err := json.Unmarshal(data, &v)
	return &v, err
}
func ReadRequestFileMessage(data []byte) (*RequestFileMessage, error) {
	var v RequestFileMessage
	error := json.Unmarshal(data, &v)
	return &v, error
}
func ReadRequestDirectoryMessage(data []byte) (*RequestDirectoryMessage, error) {
	var v RequestDirectoryMessage
	error := json.Unmarshal(data, &v)
	return &v, error
}

func ReadDirectoryStructureMessage(data []byte) (*DirectoryStructureMessage, error) {
	var v DirectoryStructureMessage
	error := json.Unmarshal(data, &v)
	return &v, error
}
