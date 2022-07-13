package connections

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"HiddenGalleryHub/server/models"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/elliotchance/pie/v2"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Pool struct {
	Upgrader           *websocket.Upgrader
	Done               chan struct{}
	Size               int
	Connections        []*websocket.Conn
	ConnectionMap      map[string]*websocket.Conn
	MutexMap           map[string]*sync.Mutex
	CurrentFileChannel map[string]chan []byte
}

func CreatePool() *Pool {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &Pool{
		Upgrader:           &upgrader,
		Done:               make(chan struct{}),
		Size:               0,
		Connections:        make([]*websocket.Conn, 0),
		ConnectionMap:      make(map[string]*websocket.Conn),
		MutexMap:           make(map[string]*sync.Mutex),
		CurrentFileChannel: make(map[string]chan []byte),
	}
}

func (pool *Pool) AddWsConnection(w http.ResponseWriter, r *http.Request, db *gorm.DB) error {
	c, err := pool.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	pool.Connections = append(pool.Connections, c)
	pool.Size += 1
	// var machine models.Machine
	// remoteAddr := strings.Split(r.RemoteAddr, ":")
	// ip := remoteAddr[0]
	// port := remoteAddr[1]
	// result := db.Where(&models.Machine{
	// 	LatestIp:   ip,
	// 	LatestPort: port,
	// }).Find(&machine)
	// if result.RowsAffected > 0 {
	// 	machine.IsOnline = true
	// 	db.UpdateColumn("is_online", &machine)
	// } else {
	// 	fmt.Println("unknown socket connected!")
	// }
	go addAndRunListeners(c, pool, db)
	return nil
}

func addAndRunListeners(conn *websocket.Conn, pool *Pool, db *gorm.DB) {
	var currentMachine models.Machine
	log.Println("sending init message")
	ws.SendMessage(conn, messages.MessageTypeInit, nil)
	for {
		mt, message, err := ws.RecvMessage(conn)
		if err != nil {
			log.Println("read message failed")
			removeConnectionFromPool(currentMachine.Name, conn, pool, db)
			return
		}
		switch mt {
		case messages.MessageTypeMachineInfo:
			{
				msg, formatError := messages.ReadMachineInfoFromMessage(message)
				if formatError != nil {
					log.Println("read machine message failed")
					removeConnectionFromPool(currentMachine.Name, conn, pool, db)
					return
				}
				currentMachine = *pool.ProcessMachineInfo(msg, conn, db)
				break
			}
		case messages.MessageTypeDirectoryStructure:
			{
				msg, _ := messages.ReadDirectoryStructureMessage(message)
				ProcessDirectoryStructure(&currentMachine, msg, db)
				break
			}
		case messages.MessageTypeFileContent:
			{
				pool.CurrentFileChannel[currentMachine.Name] <- message
				break
			}
		case messages.MessageTypeFileFinish:
			{
				close(pool.CurrentFileChannel[currentMachine.Name])
				break
			}
		case messages.MessageTypeFileException:
			{
				close(pool.CurrentFileChannel[currentMachine.Name])
				break
			}
		}
	}

}

func removeConnectionFromPool(name string, conn *websocket.Conn, pool *Pool, db *gorm.DB) {
	originalSize := pool.Size
	pool.Connections = pie.FilterNot(pool.Connections, func(c *websocket.Conn) bool { return c == conn })
	pool.ConnectionMap[name] = nil
	pool.Size = len(pool.Connections)
	if originalSize != pool.Size+1 {
		panic("a connection not in pool is removed")
	}
	fmt.Println("a connection is removed from pool")
	var machine models.Machine
	result := db.Where(models.Machine{
		Name: name,
	}).First(&machine)
	if result.RowsAffected > 0 {
		machine.IsOnline = false
		db.Save(machine)
	}
}
