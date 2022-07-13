package connections

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"HiddenGalleryHub/server/models"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func (pool *Pool) ProcessMachineInfo(msg *messages.MachineInfoMessage, conn *websocket.Conn, db *gorm.DB) *models.Machine {
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
	ip := remoteAddr[0]
	port := remoteAddr[1]
	var machine models.Machine
	result := db.Where(models.Machine{
		Name: msg.MachineName,
	}).First(&machine)
	if result.RowsAffected > 0 {
		if machine.PasswdSum != "" && machine.PasswdSum != msg.PasswdSum {
			conn.Close()
			return nil
		}
		machine.IsOnline = true
		machine.LatestIp = ip
		machine.LatestPort = port
		machine.Name = msg.MachineName
		machine.PasswdSum = msg.PasswdSum
		db.Save(machine)
	} else {
		machine = models.Machine{
			LatestIp:   ip,
			LatestPort: port,
			Name:       msg.MachineName,
			IsOnline:   true,
			PasswdSum:  msg.PasswdSum,
		}
		db.Create(&machine)
	}
	dirReq := messages.RequestDirectoryMessage{
		RelativePath: ".",
	}
	dirReqData, err := json.Marshal(&dirReq)
	if err != nil {
		log.Fatal(err)
	}

	pool.ConnectionMap[machine.Name] = conn
	pool.MutexMap[machine.Name] = &sync.Mutex{}
	// get directory structure after machine is online.
	err = ws.SendMessage(conn, messages.MessageTypeRequestDirectory, dirReqData)
	if err != nil {
		log.Fatal(err)
	}
	return &machine
}

func ProcessDirectoryStructure(machine *models.Machine, msg *messages.DirectoryStructureMessage, db *gorm.DB) {
	// set all entries in the machine as invalid first
	db.Model(&models.Directory{}).Where(&models.Directory{
		MachineId: machine.ID,
	}).UpdateColumn("is_invalid", 1)
	db.Model(&models.FileEntry{}).Where(&models.FileEntry{
		MachineId: machine.ID,
	}).UpdateColumn("is_invalid", 1)

	dirIdMap := make(map[string]uint)

	// firstly update directories
	for _, dir := range msg.DirectoryEntries {
		var readingDirectory models.Directory
		result := db.Where(models.Directory{
			RelativePath: dir.RelativePath,
			MachineId:    machine.ID,
		}).Find(&readingDirectory)
		if result.RowsAffected > 0 {
			readingDirectory.IsInvalid = false
			db.Save(&readingDirectory)
		} else {
			readingDirectory = models.Directory{
				Name:              dir.Name,
				RelativePath:      dir.RelativePath,
				MachineId:         machine.ID,
				IsInvalid:         false,
				IsRootDirectory:   dir.RelativePath == ".",
				ParentDirectoryId: dirIdMap[dir.ParentRelativePath],
			}
			db.Create(&readingDirectory)
		}
		dirIdMap[readingDirectory.RelativePath] = readingDirectory.ID
	}
	// then update files
	for _, file := range msg.FileEntries {
		var readingFileEntry models.FileEntry
		result := db.Where(models.FileEntry{
			RelativePath: file.RelativePath,
			MachineId:    machine.ID,
		}).Find(&readingFileEntry)
		if result.RowsAffected > 0 {
			readingFileEntry.IsInvalid = false
			db.Save(readingFileEntry)
		} else {
			readingFileEntry = models.FileEntry{
				Name:              file.Name,
				RelativePath:      file.RelativePath,
				MachineId:         machine.ID,
				IsInvalid:         false,
				ParentDirectoryId: dirIdMap[file.ParentRelativePath],
			}
			db.Create(&readingFileEntry)
		}
	}
}
