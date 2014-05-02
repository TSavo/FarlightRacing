package greatspacerace

import "github.com/TSavo/chipmunk/vect"

type Wall struct {
	Point1, Point2 vect.Vect
}

type Checkpoint struct {
	Wall                        Wall
	Angle, XCrossing, YCrossing vect.Float
}

type Track struct {
	Id, Name    string
	Walls       []Wall
	GoalLine    Checkpoint
	Checkpoints []Checkpoint
}

func (this *Checkpoint) GetStartingPositions(pieces int) []vect.Vect {
	xDif := (this.Wall.Point2.X - this.Wall.Point1.X) / vect.Float(pieces+1)
	yDif := (this.Wall.Point2.Y - this.Wall.Point1.Y) / vect.Float(pieces+1)
	places := make([]vect.Vect, pieces)
	for x := 0; x < pieces; x++ {
		places[x].X = this.Wall.Point1.X + (xDif * vect.Float(x+1))
		places[x].Y = this.Wall.Point1.Y + (yDif * vect.Float(x+1))
	}
	return places
}
