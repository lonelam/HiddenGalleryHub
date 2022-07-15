package processor

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (c *WsClientConnection) onRequestFileSendFilePieceByPiece(message []byte) {
	c.fileMu.Lock()
	defer c.fileMu.Unlock()
	requestFile, err := messages.ReadRequestFileMessage(message)
	if err != nil {
		log.Println(err)
		c.conn.Close()
		return
	}
	openFileHandler, err := os.OpenFile(filepath.Join(c.rootDir, requestFile.RelativePath), os.O_RDONLY, 0755)
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
		sendingBuffer := buffer[:currentReadSize]
		log.Printf("file buffer read(%d) send with len(%d)\n", currentReadSize, len(sendingBuffer))
		err = ws.SendMessage(c.conn, messages.MessageTypeFileContent, sendingBuffer)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
