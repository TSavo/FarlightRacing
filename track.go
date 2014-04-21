package slotserver

type Track struct {
	Id, Name      string
	Pieces        []Piece
	Lanes         []Lane
	StartingPoint StartingPoint
}

type Piece struct {
	Length, Radius, Angle float64
	Switch                bool
}

type Lane struct {
	DistanceFromCenter float64
	Index              int
}

type StartingPoint struct {
	Position Point
	Angle    float64
}

type Point struct {
	X, Y float64
}
