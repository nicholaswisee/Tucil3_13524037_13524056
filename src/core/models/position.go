package models

import "fmt"

type Position struct {
	X, Y int
}

func (p Position) Add(dx, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}

func (p Position) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}
