package greatspacerace

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/TSavo/chipmunk"
	"github.com/TSavo/chipmunk/vect"
	"github.com/twinj/uuid"
	"math"
)

type Race struct {
	Id      string
	Track   *Track
	Ships   []*Ship
	Space   *chipmunk.Space
	Started bool
}

func NewRace(track *Track) *Race {
	race := &Race{uuid.NewV4().String(), track, make([]*Ship, 0), chipmunk.NewSpace(), false}
	staticBody := chipmunk.NewBodyStatic()
	for _, wall := range track.Walls {
		segment := chipmunk.NewSegment(wall.Point1, wall.Point2, 0)
		staticBody.AddShape(segment)
	}
	race.Space.AddBody(staticBody)
	return race
}

type PlayerPosition struct {
	Name                 string
	Dimensions, Position vect.Vect
	Angle                vect.Float
}

func (this *Race) RunRace() {
	this.StartRace()
	for {
		players := make([]PlayerPosition, len(this.Ships))
		for x, ship := range this.Ships {
			players[x] = PlayerPosition{ship.Player.Name, vect.Vect{ship.Prototype.Width, ship.Prototype.Height}, ship.Body.Position(), ship.Body.Angle()}
		}
		for _, ship := range this.Ships {
			ship.Player.Send("RaceUpdate", players)
		}
		for _, ship := range this.Ships {
			reader := bufio.NewReader(ship.Player.Conn)
			message, _, _ := reader.ReadLine()
			decoded := new(map[string]float64)
			json.Unmarshal(message, &decoded)
			ship.Controller.Thrust = vect.Float((*decoded)["Thrust"])
			ship.Controller.Turning = vect.Float((*decoded)["Rotation"])
			ship.ApplyThrust(ship.Controller.Thrust)
			ship.ApplyRotation(ship.Controller.Turning)
		}
		this.MoveShips()
	}
}

func (this *Race) MoveShips() {
	old := make(map[string]vect.Vect)
	for _, ship := range this.Ships {
		old[ship.Player.Name] = ship.Body.Position()
	}
	this.StepRace(1.0 / 60.0)
//	for _, ship := range this.Ships {
//		//before := old[ship.Player.Name]
//	}

}

func (this *Race) StartRace() {
	this.Started = true
	startPoints := this.Track.GoalLine.GetStartingPositions(len(this.Ships))
	for x, ship := range this.Ships {
		ship.Body.SetPosition(startPoints[x])
		ship.Body.SetAngle(vect.Float(this.Track.GoalLine.Angle * 2 * math.Pi))
		this.Space.AddBody(ship.Body)
		fmt.Println("added", ship.Body)
	}
}

func (this *Race) RegisterRacer(player *Player, prototype *Prototype) *Ship {
	box := chipmunk.NewBox(vect.Vector_Zero, prototype.Width, prototype.Height)
	box.SetElasticity(0.9)
	body := chipmunk.NewBody(prototype.Mass, box.Moment(prototype.Moment))
	body.AddShape(box)
	ship := NewShip(player, prototype, body)
	this.Ships = append(this.Ships, ship)
	return ship
}

func (this *Race) StepRace(dt vect.Float) {
	this.Space.Step(dt)
}
