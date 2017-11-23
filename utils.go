package gwk

import (
	"github.com/Luncher/gwk/pkg/structs"
)

func isPointInRect(point *structs.Point, rect *structs.Rect) bool {
	return point.X >= rect.X && point.Y >= rect.Y && point.X < (rect.X+rect.W) && point.Y < (rect.Y+rect.H)
}
