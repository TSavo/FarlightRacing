// +build ignore

package main

import (
	"encoding/json"
	. "github.com/TSavo/GreatSpaceRace"
	"github.com/TSavo/chipmunk/vect"
	"net"
	//"fmt"
	"bufio"
)

type PlayerPosition struct {
	Name                 string
	Dimensions, Position vect.Vect
	Angle                vect.Float
}

type Server struct {
	Races []*Race
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
	}
	conn, err := ln.Accept()
	if err != nil {
		return
	}

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

	track := &Track{"testtrack", "Test Track", walls, goal}
	player := Player{Name: "Test Player", Conn: conn}
	race := NewRace(track)
	prototype := Prototype{"Test Ship", 100, 50, 200000, 3500000, 1000, 1000}
	race.RegisterRacer(&player, &prototype)
	race.StartRace()
	for {
		message := make(map[string]interface{})
		message["track"] = track
		players := make([]PlayerPosition, len(race.Ships))
		for x, ship := range race.Ships {
			players[x] = PlayerPosition{ship.Player.Name, vect.Vect{ship.Prototype.Width, ship.Prototype.Height}, ship.Body.Position(), ship.Body.Angle()}
		}
		message["players"] = players
		for _, ship := range race.Ships {
			reader := bufio.NewReader(conn)
			writer := json.NewEncoder(conn)
			writer.Encode(message)
			message, _, _ := reader.ReadLine()
			decoded := new(map[string]float64)
			json.Unmarshal(message, &decoded)
			ship.Controller.Thrust = vect.Float((*decoded)["Thrust"])
			ship.Controller.Turning = vect.Float((*decoded)["Rotation"])
			ship.ApplyThrust(ship.Controller.Thrust)
			ship.ApplyRotation(ship.Controller.Turning)
		}
		race.StepRace(1.0 / 60.0)
	}
}
