package greatspacerace

import (
	"fmt"
	"github.com/TSavo/chipmunk"
	"github.com/TSavo/chipmunk/vect"
	"math"
)

type Race struct {
	Track   *Track
	Ships   []*Ship
	Space   *chipmunk.Space
	Started bool
}

func NewRace(track *Track) *Race {
	race := &Race{track, make([]*Ship, 0), chipmunk.NewSpace(), false}
	staticBody := chipmunk.NewBodyStatic()
	for _, wall := range track.Walls {
		segment := chipmunk.NewSegment(wall.Point1, wall.Point2, 0)
		staticBody.AddShape(segment)
	}
	race.Space.AddBody(staticBody)
	return race
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
