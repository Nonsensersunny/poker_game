package client

type Position struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

func NewPosition(x1, y1, x2, y2 int) Position {
	return Position{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}
}
