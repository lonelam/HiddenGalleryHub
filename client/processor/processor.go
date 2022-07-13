package processor

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

func (c *WsClientConnection) onInitSendMachineInfo() {
	machineInfo := messages.MachineInfoMessage{
		MachineName: "9527",
		PasswdSum:   "",
	}
	machineInfoMsg, _ := json.Marshal(machineInfo)
	ws.SendMessage(c.conn, messages.MessageTypeMachineInfo, machineInfoMsg)
}

func (c *WsClientConnection) onRequestFileSendFilePieceByPiece(message []byte) {

	requestFile, err := messages.ReadRequestFileMessage(message)
	if err != nil {
		log.Fatal(err)
		return
	}
	openFileHandler, err := os.OpenFile(path.Join(c.rootDir, requestFile.RelativePath), os.O_RDONLY, 0755)
	if err != nil {
		log.Printf("open file failed. err: %v.\n ", err)
		ws.SendMessage(c.conn, messages.MessageTypeFileException, nil)
		return
	}
	ws.SendMessage(c.conn, messages.MessageTypeFileStart, nil)
	defer ws.SendMessage(c.conn, messages.MessageTypeFileFinish, nil)
	buffer := make([]byte, c.config.ClientBufferSize)
	for {
		currentReadSize, err := openFileHandler.Read(buffer)
		if err == io.EOF && currentReadSize == 0 {
			return
		}
		if err != nil {
			log.Printf("reading file error: %v", err)
			break
		}
		err = ws.SendMessage(c.conn, messages.MessageTypeFileContent, buffer)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func (c *WsClientConnection) onRequestDirectorySendAllDirectories(message []byte) {
	requestDirectory, err := messages.ReadRequestDirectoryMessage(message)
	if err != nil {
		log.Fatal(err)
		return
	}
	startPath := path.Join(c.rootDir, requestDirectory.RelativePath)
	startRelPath, _ := filepath.Rel(c.rootDir, startPath)

	directoryArr := make([]messages.DirectoryEntry, 1)
	directoryArr[0] = messages.DirectoryEntry{
		Name:               path.Base(startPath),
		RelativePath:       startRelPath,
		ParentRelativePath: path.Dir(startRelPath),
	}
	fileArr := make([]messages.FileEntry, 0)
	// log.Printf("root:%s relative: %s startpath: %s\n", c.rootDir, requestDirectory.RelativePath, startPath)
	collectDirectoryStructure(c.rootDir, startPath, &directoryArr, &fileArr)
	dirStruct := messages.DirectoryStructureMessage{
		DirectoryEntries: directoryArr,
		FileEntries:      fileArr,
	}
	dirStructMsg, _ := json.Marshal(dirStruct)
	ws.SendMessage(c.conn, messages.MessageTypeDirectoryStructure, dirStructMsg)
}

func collectDirectoryStructure(rootPath string, startPath string, directoryArr *[]messages.DirectoryEntry, fileArr *[]messages.FileEntry) {
	parentRelPath, _ := filepath.Rel(rootPath, startPath)
	dirEntries, err := os.ReadDir(startPath)
	if err != nil {
		log.Fatalf("readDir %s failed, err: %v", startPath, err)
		return
	}
	log.Printf("parsing %s\n", parentRelPath)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			*directoryArr = append(*directoryArr, messages.DirectoryEntry{
				Name:               dirEntry.Name(),
				RelativePath:       path.Join(parentRelPath, dirEntry.Name()),
				ParentRelativePath: parentRelPath,
			})
			collectDirectoryStructure(rootPath, path.Join(startPath, dirEntry.Name()), directoryArr, fileArr)
		} else {
			*fileArr = append(*fileArr, messages.FileEntry{
				Name:               dirEntry.Name(),
				RelativePath:       path.Join(parentRelPath, dirEntry.Name()),
				ParentRelativePath: parentRelPath,
			})
		}
	}
}
