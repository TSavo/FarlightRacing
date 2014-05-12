package greatspacerace

import "github.com/TSavo/chipmunk/vect"

type Segment struct {
	Point1, Point2 vect.Vect
}

func (S1 *Segment) Intersects(S2 *Segment) bool {
	u := vect.Sub(S1.Point2, S1.Point1)
	v := vect.Sub(S2.Point2, S2.Point1)
	w := vect.Sub(S1.Point1, S2.Point1)
	D := vect.Cross(u, v)

	// test if  they are parallel (includes either being a point)
	if vect.FAbs(D) < 0.0000001 { // S1 and S2 are parallel
		if vect.Cross(u, w) != 0 || vect.Cross(v, w) != 0 {
			return false // they are NOT collinear
		}
		// they are collinear or degenerate
		// check if they are degenerate  points
		du := vect.Dot(u, u)
		dv := vect.Dot(v, v)
		if du == 0 && dv == 0 { // both segments are points
			if !vect.Equals(S1.Point1, S2.Point1) { // they are distinct  points
				return false
			}
			return true
		}
		if du == 0 { // S1 is a single point
			if !S2.Contains(S1.Point1) { // but is not in S2
				return false
			}
			return true
		}
		if dv == 0 { // S2 a single point
			if !S1.Contains(S2.Point1) { // but is not in S1
				return false
			}
			return true
		}
		// they are collinear segments - get  overlap (or not)
		var t0, t1 vect.Float // endpoints of S1 in eqn for S2
		w2 := vect.Sub(S1.Point2, S2.Point1)
		if v.X != 0 {
			t0 = w.X / v.X
			t1 = w2.X / v.X
		} else {
			t0 = w.Y / v.Y
			t1 = w2.Y / v.Y
		}
		if t0 > t1 { // must have t0 smaller than t1
			t0, t1 = t1, t0 // swap if not
		}
		if t0 > 1 || t1 < 0 {
			return false // NO overlap
		}
		return true
	}

	// the segments are skew and may intersect in a point
	// get the intersect parameter for S1
	sI := vect.Cross(v, w) / D
	if sI < 0 || sI > 1 { // no intersect with S1
		return false
	}
	// get the intersect parameter for S2
	tI := vect.Cross(u, w) / D
	if tI < 0 || tI > 1 { // no intersect with S2
		return false
	}
	return true
}

func (S *Segment) Contains(P vect.Vect) bool {
	if S.Point1.X != S.Point2.X { // S is not  vertical
		if S.Point1.X <= P.X && P.X <= S.Point2.X {
			return true
		}
		if S.Point1.X >= P.X && P.X >= S.Point2.X {
			return true
		}
	} else { // S is vertical, so test y  coordinate
		if S.Point1.Y <= P.Y && P.Y <= S.Point2.Y {
			return true
		}
		if S.Point1.Y >= P.Y && P.Y >= S.Point2.Y {
			return true
		}
	}
	return false
}

type Track struct {
	Id, Name      string
	Laps          int
	Segments      []Segment
	Goal          Segment
	StartingAngle vect.Float
	Checkpoints   []Segment
}

func (this *Segment) GetStartingPositions(pieces int) []vect.Vect {
	xDif := (this.Point2.X - this.Point1.X) / vect.Float(pieces+1)
	yDif := (this.Point2.Y - this.Point1.Y) / vect.Float(pieces+1)
	places := make([]vect.Vect, pieces)
	for x := 0; x < pieces; x++ {
		places[x].X = this.Point1.X + (xDif * vect.Float(x+1))
		places[x].Y = this.Point1.Y + (yDif * vect.Float(x+1))
	}
	return places
}
