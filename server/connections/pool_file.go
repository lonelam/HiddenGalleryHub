package connections

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"HiddenGalleryHub/server/models"
	"encoding/json"
	"errors"
	"log"
)

func (pool *Pool) GetFileFromRemote(fileEntry *models.FileEntry) (chan []byte, error) {
	resultChan := make(chan []byte)
	if !fileEntry.Machine.IsOnline {
		return nil, errors.New("machine is offline")
	}
	go func() {
		defer close(resultChan)
		bufferChan := make(chan []byte)
		mu := pool.MutexMap[fileEntry.Machine.Name]
		conn := pool.ConnectionMap[fileEntry.Machine.Name]
		mu.Lock()
		defer mu.Unlock()
		pool.CurrentFileChannel[fileEntry.Machine.Name] = bufferChan
		defer func() { pool.CurrentFileChannel[fileEntry.Machine.Name] = nil }()
		reqMsg, _ := json.Marshal(&messages.RequestFileMessage{
			RelativePath: fileEntry.RelativePath,
		})
		ws.SendMessage(conn, messages.MessageTypeRequestFile, reqMsg)
		for {
			buffer, ok := <-bufferChan
			log.Printf("buffer received: %d\n", len(buffer))
			if !ok {
				log.Printf("buffer closed.\n")
				break
			}
			resultChan <- buffer
		}
	}()
	return resultChan, nil
}
