package gwk

type Position struct {
	x int
	y int
}

func NewPosition(x, y int) *Position {
	return &Position{x, y}
}

type Rect struct {
	x, y, w, h int
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{x, y, w, h}
}

type Point struct {
	x, y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}
