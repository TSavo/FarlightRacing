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
	RegisterChan chan *JoinMessage
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
	walls := []Wall{
		Wall{vect.Vect{0, 0}, vect.Vect{0, 900}},
		Wall{vect.Vect{0, 900}, vect.Vect{1500, 900}},
		Wall{vect.Vect{1500, 900}, vect.Vect{1500, 0}},
		Wall{vect.Vect{1500, 0}, vect.Vect{0, 0}},

		Wall{vect.Vect{250, 250}, vect.Vect{250, 750}},
		Wall{vect.Vect{250, 750}, vect.Vect{750, 750}},
		Wall{vect.Vect{750, 750}, vect.Vect{750, 250}},
		Wall{vect.Vect{750, 250}, vect.Vect{250, 250}},
	}

	goalWall := Wall{vect.Vect{500, 0}, vect.Vect{500, 250}}

	goal := GoalLine{goalWall, 0, 1, 0}

	return &Track{"testtrack", "Test Track", walls, goal}
}

func (this *Server) RegisterThread() {
	for {
		joinMessage := <-this.RegisterChan
		player := Player{Name: joinMessage.Name, Conn: joinMessage.Conn}
		var lobby *RaceLobby = nil
		for _, race := range this.Lobbies {
			if race.Race.Track.Id == joinMessage.Track && race.NumPlayers == joinMessage.NumPlayers && race.Password == joinMessage.Password {
				lobby = race
				break
			}
		}
		if(lobby == nil) {
			lobby = &RaceLobby{NewRace(GetTrack(joinMessage.Track)), joinMessage.NumPlayers, joinMessage.Password}
			this.Lobbies = append(this.Lobbies, lobby)
		}
		lobby.Race.RegisterRacer(&player, GetPrototype(joinMessage.Prototype))
		trackMessage := make(map[string]interface{})
		trackMessage["Track"] = lobby.Race.Track
		trackMessage["Race"] = lobby.Race.Id
		player.Send("RaceData", trackMessage)
		this.StartRaces();
	}
}

func (this *Server) StartRaces(){
	for x := 0; x < len(this.Lobbies); x++ {
		race := this.Lobbies[x]
		if race.NumPlayers == len(race.Race.Ships) {
			this.Races = append(this.Races, race.Race)
			this.Lobbies = append(this.Lobbies[:x], this.Lobbies[x+1:]...)
			x--
			go race.Race.RunRace()
		}
	}
}

func (this *Server) HandleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	message, _, _ := reader.ReadLine()
	joinMessage := JoinMessage{}
	json.Unmarshal(message, &joinMessage)
	joinMessage.Conn = conn
	this.RegisterChan <- &joinMessage
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
		go this.HandleConnection(conn)
	}
}
