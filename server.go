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

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		walls := []Wall{
			Wall{vect.Vect{0, 0}, vect.Vect{0, 1000}},
			Wall{vect.Vect{0, 1000}, vect.Vect{1000, 1000}},
			Wall{vect.Vect{1000, 1000}, vect.Vect{1000, 0}},
			Wall{vect.Vect{1000, 0}, vect.Vect{0, 0}},

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
		prototype := Prototype{"Test Ship", 100, 50, 1000000, 1000000000, 1000, 1000}
		race.RegisterRacer(&player, &prototype)
		encoder := json.NewEncoder(conn)
		//decoder := json.Decoder(conn)
		race.StartRace()
		for {
			message := make(map[string]interface{})
			message["track"] = track
			players := make([]PlayerPosition, len(race.Ships))
			for x, ship := range race.Ships {
				reader := bufio.NewReader(conn)
				message, _, _ := reader.ReadLine()
				//fmt.Println(string(message))
				decoded := new(map[string]float64)
				json.Unmarshal(message, &decoded)
				//
				//fmt.Println(*decoded)
				players[x] = PlayerPosition{ship.Player.Name, vect.Vect{ship.Prototype.Width, ship.Prototype.Height}, ship.Body.Position(), ship.Body.Angle()}
				ship.Controller.Thrust = vect.Float((*decoded)["Thrust"])
				ship.Controller.Turning = vect.Float((*decoded)["Rotation"])
				ship.ApplyThrust(ship.Controller.Thrust)
				ship.ApplyRotation(ship.Controller.Turning)
			}
			message["players"] = players
			encoder.Encode(message)
			race.StepRace(1.0 / 60.0)
		}
	}
	// set up physics
}
