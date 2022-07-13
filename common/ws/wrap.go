package ws

import (
	"HiddenGalleryHub/common/pb"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func SendMessage(conn *websocket.Conn, messageType int, data []byte) error {
	msg := pb.MsgWrapper{
		MsgType: int32(messageType),
		Content: data,
	}
	byteData, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	// log.Printf("Message send. mt: %d", messageType)
	err = conn.WriteMessage(websocket.BinaryMessage, byteData)

	if err != nil {
		return err
	}
	return nil
}
func RecvMessage(conn *websocket.Conn) (int, []byte, error) {
	log.Printf("start listening next message \n")
	mt, message, err := conn.ReadMessage()
	if err != nil {
		return -1, nil, err
	}
	if mt != websocket.BinaryMessage {
		log.Fatalf("message type except binary received")
		return -1, nil, err
	}
	var wrappedMsg pb.MsgWrapper
	err = proto.Unmarshal(message, &wrappedMsg)
	// log.Printf("read message success, mt: %d\n", wrappedMsg.GetMsgType())
	if err != nil {
		return -1, nil, err
	}

	return int(wrappedMsg.GetMsgType()), wrappedMsg.GetContent(), err
}
