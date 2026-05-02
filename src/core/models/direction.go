package models

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var Directions = []Direction{Up, Right, Down, Left}

func (d Direction) DxDy() (dx, dy int) {
	switch d {
	case Up:
		return -1, 0
	case Right:
		return 0, 1
	case Down:
		return 1, 0
	case Left:
		return 0, -1
	}
	return 0, 0
}

func (d Direction) DirectionName() string {
	switch d {
	case Up:
		return "U"
	case Right:
		return "R"
	case Down:
		return "D"
	case Left:
		return "L"
	}
	return "?"
}
