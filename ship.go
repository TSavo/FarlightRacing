package greatspacerace

import (
	"github.com/TSavo/chipmunk"
	"github.com/TSavo/chipmunk/vect"
	//"fmt"
)

type Ship struct {
	Player         *Player
	Prototype      *Prototype
	Body           *chipmunk.Body
	Controller     *Controller
	NextCheckpoint int
	Laps           []Lap
}

type Lap struct {
	Start, Finish int64
}

type Prototype struct {
	Name                                         string
	Width, Height, Thrust, Turning, Mass, Moment vect.Float
}

type Controller struct {
	Thrust, Turning vect.Float
}

func NewShip(player *Player, prototype *Prototype, body *chipmunk.Body) *Ship {
	return &Ship{player, prototype, body, &Controller{0, 0}, 0, make([]Lap, 0)}
}

func (this *Ship) ApplyThrust(thrust vect.Float) {
	thrustVector := vect.FromAngle(this.Body.Angle())
	thrustVector.Normalize()
	thrustVector.Mult(this.Prototype.Thrust * thrust)
	this.Body.AddForce(thrustVector.X, thrustVector.Y)
}

func (this *Ship) ApplyRotation(thrust vect.Float) {
	this.Body.SetTorque(this.Prototype.Turning * thrust)
}
