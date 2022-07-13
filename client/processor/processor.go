package processor

import (
	"HiddenGalleryHub/common/messages"
	"HiddenGalleryHub/common/ws"
	"encoding/json"
)

func (c *WsClientConnection) onInitSendMachineInfo() {
	machineInfo := messages.MachineInfoMessage{
		MachineName: "9527",
		PasswdSum:   "",
	}
	machineInfoMsg, _ := json.Marshal(machineInfo)
	ws.SendMessage(c.conn, messages.MessageTypeMachineInfo, machineInfoMsg)
}
