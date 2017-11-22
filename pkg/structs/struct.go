package structs

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) *Position {
	return &Position{x, y}
}

type Rect struct {
	X, Y, W, H int
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{x, y, w, h}
}

type Point struct {
	X, Y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}
