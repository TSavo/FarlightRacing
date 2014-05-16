package greatspacerace

import (
	"encoding/json"
	"github.com/TSavo/chipmunk/vect"
	"net"
	//"fmt"
	"bufio"
	"strconv"
)

type RaceLobby struct {
	*Race
	NumPlayers int
	Password   string
}

type Server struct {
	Races        []*Race
	Lobbies      []*RaceLobby
	RegisterChan chan *net.Conn
}

func NewServer() *Server {
	return &Server{make([]*Race, 0), make([]*RaceLobby, 0), make(chan *net.Conn, 100)}
}

type JoinMessage struct {
	Name, Track, Prototype, Password string
	NumPlayers                       int
	Conn                             net.Conn
}

func GetPrototype(name string) *Prototype {
	return &Prototype{"Test Ship", 100, 50, 200000, 3500000, 1000, 1000}
}

func GetTrack(name string) *Track {
	Segments := []Segment{
		Segment{vect.Vect{0, 0}, vect.Vect{0, 900}},
		Segment{vect.Vect{0, 900}, vect.Vect{1500, 900}},
		Segment{vect.Vect{1500, 900}, vect.Vect{1500, 0}},
		Segment{vect.Vect{1500, 0}, vect.Vect{0, 0}},

		Segment{vect.Vect{250, 250}, vect.Vect{250, 750}},
		Segment{vect.Vect{250, 750}, vect.Vect{750, 750}},
		Segment{vect.Vect{750, 750}, vect.Vect{750, 250}},
		Segment{vect.Vect{750, 250}, vect.Vect{250, 250}},
	}

	goalSegment := Segment{vect.Vect{500, 0}, vect.Vect{500, 250}}


	checkpoints := []Segment{
		Segment{vect.Vect{750, 450}, vect.Vect{1500, 450}},
		Segment{vect.Vect{500, 750}, vect.Vect{500, 900}},
		Segment{vect.Vect{0, 450}, vect.Vect{250, 450}},
	}

	return &Track{name, name, 3, 25000, Segments, goalSegment, 0, checkpoints}
}

func (this *Server) RegisterThread() {
	for {
		conn := <-this.RegisterChan
		reader := bufio.NewReader(*conn)
		message, _, _ := reader.ReadLine()
		joinMessage := JoinMessage{}
		json.Unmarshal(message, &joinMessage)
		joinMessage.Conn = *conn
		ship := NewShip(GetPrototype(joinMessage.Prototype))
		player := NewPlayer(joinMessage.Name, &joinMessage.Conn, ship)
		var lobby *RaceLobby = nil
		for _, race := range this.Lobbies {
			if race.Race.Track.Id == joinMessage.Track && race.NumPlayers == joinMessage.NumPlayers && race.Password == joinMessage.Password {
				lobby = race
				break
			}
		}
		if lobby == nil {
			lobby = &RaceLobby{NewRace(GetTrack(joinMessage.Track)), joinMessage.NumPlayers, joinMessage.Password}
			this.Lobbies = append(this.Lobbies, lobby)
		}
		lobby.Race.RegisterRacer(player)
		trackMessage := make(map[string]interface{})
		trackMessage["Track"] = lobby.Race.Track
		trackMessage["Race"] = lobby.Race.Id
		player.SendChan <- &Message{"RaceData", trackMessage, -1}
		this.StartRaces()
	}
}

func (this *Server) StartRaces() {
	for x := 0; x < len(this.Lobbies); x++ {
		lobby := this.Lobbies[x]
		if lobby.NumPlayers == len(lobby.Race.Racers) {
			this.Races = append(this.Races, lobby.Race)
			this.Lobbies = append(this.Lobbies[:x], this.Lobbies[x+1:]...)
			x--
			go lobby.Race.RunRace()
		}
	}
}

func (this *Server) Listen(port int) {
	go this.RegisterThread()
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		this.RegisterChan <- &conn
	}
}
