package greatspacerace

import (
	"encoding/json"
	"net"
)

type Player struct {
	Name string
	Conn net.Conn
	Ship Ship
}

func (this *Player) SendRaceStart() {
	this.Send("RaceStatus", "START")
}

func (this *Player) Send(messageType string, message interface{}) {
	outMessage := make(map[string]interface{})
	outMessage["Type"] = messageType
	outMessage["Data"] = message
	json.NewEncoder(this.Conn).Encode(outMessage)
}
