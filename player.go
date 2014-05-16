package greatspacerace

import (
	"bufio"
	"encoding/json"
	"net"
	"time"
)

type Player struct {
	Name        string
	Conn        *net.Conn
	Ship        *Ship
	SendChan    chan<- *Message
	ReceiveChan <-chan *Message
}

type Message struct {
	Type string
	Data interface{}
	Tick int
}

func NewPlayer(name string, conn *net.Conn, ship *Ship) *Player {
	p := Player{Name: name, Conn: conn, Ship: ship}
	p.Init()
	return &p
}

func (this *Player) Init() {
	sendChan := make(chan *Message, 100)
	this.SendChan = sendChan
	go func() {
		encoder := json.NewEncoder(*this.Conn)
		for {
			m := <-sendChan
			encoder.Encode(*m)
			if m.Type == "RaceFinished" {
				go func() {
					time.Sleep(time.Second)
					(*this.Conn).Close()
				}()
				return
			}
		}
	}()
	receiveChan := make(chan *Message, 100)
	this.ReceiveChan = receiveChan
	go func() {
		scanner := bufio.NewScanner(*this.Conn)
		defer func() {
			recover()
		}()
		for scanner.Scan() {
			m := Message{}
			json.Unmarshal(scanner.Bytes(), &m)
			receiveChan <- &m
		}
	}()
}

func (this *Player) SendRaceStart() {
	this.Send(&Message{"RaceStatus", "START", -1})
}

func (this *Player) Send(message *Message) {
	this.SendChan <- message
}
