package processor

import (
	"HiddenGalleryHub/client/constants"
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsClientConnection struct {
	conn    *websocket.Conn
	rootDir string
	config  *constants.AppConfiguration
	fileMu  *sync.Mutex
	name    string
}

func CreateWsConnection(url string, rootDir string, name string) (*WsClientConnection, error) {

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("dial:", err)
		return nil, err
	}
	return &WsClientConnection{
		conn:    c,
		rootDir: rootDir,
		config: &constants.AppConfiguration{
			ClientBufferSize: 1024 * 1024,
		},
		fileMu: &sync.Mutex{},
		name:   name,
	}, nil
}

func (c *WsClientConnection) StartUp() chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			mt, message, err := ws.RecvMessage(c.conn)
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("Message received. mt: %d\n", mt)
			switch mt {
			case messages.MessageTypeInit:
				{
					c.onInitSendMachineInfo()
					break
				}
			case messages.MessageTypeRequestDirectory:
				{
					c.onRequestDirectorySendAllDirectories(message)
					break
				}
			case messages.MessageTypeRequestFile:
				{

					c.onRequestFileSendFilePieceByPiece(message)
					break
				}
			}
		}
	}()
	return done
}
