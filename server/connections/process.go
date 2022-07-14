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

	// set all data invalid
	db.Model(&models.Directory{}).Where(&models.Directory{
		MachineId: machine.ID,
	}).UpdateColumn("is_invalid", 1)
	db.Model(&models.FileEntry{}).Where(&models.FileEntry{
		MachineId: machine.ID,
	}).UpdateColumn("is_invalid", 1)
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
			readingDirectory.Name = dir.Name
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
	tx := db.Session(&gorm.Session{
		SkipDefaultTransaction: false,
	})
	tx.Exec("BEGIN TRANSACTION;")
	for _, file := range msg.FileEntries {
		err := tx.Exec(`
		INSERT INTO
		file_entries 
		(
			created_at,
			updated_at,
			name,
			relative_path,
			file_size,
			thumbnail,
			machine_id,
			is_invalid,
			parent_directory_id,
			thumbnail_height,
			thumbnail_width
		)
		VALUES
		(
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			?,
			?,
			?,
			?,
			?,
			0,
			(SELECT id FROM directories WHERE relative_path=? and machine_id=? LIMIT 1),
			?,
			?
		)
		ON CONFLICT(relative_path, machine_id) DO
		UPDATE
		SET updated_at = EXCLUDED.updated_at,
			name = EXCLUDED.name,
			file_size = EXCLUDED.file_size,
			thumbnail = EXCLUDED.thumbnail,
			is_invalid = 0,
			parent_directory_id = EXCLUDED.parent_directory_id,
			thumbnail_height = EXCLUDED.thumbnail_height,
			thumbnail_width = EXCLUDED.thumbnail_width
		;
		`, file.Name,
			file.RelativePath,
			file.FileSize,
			file.Thumbnail,
			machine.ID,
			file.ParentRelativePath,
			machine.ID,
			file.ThumbnailHeight,
			file.ThumbnailWidth,
		).Error
		if err != nil {
			log.Printf("exec update file SQL err: %v\n", err)
		}
	}
	tx.Exec("COMMIT;")
}
