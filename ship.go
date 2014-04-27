package slotserver

import "github.com/TSavo/chipmunk"

type Ship struct {
	Player Player
	Shape *chipmunk.Shape
}