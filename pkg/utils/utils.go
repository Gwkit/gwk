package utils

import (
	"honnef.co/go/js/dom"
	"math"
)

const (
	RoundRectTL  = 1
	RoundRectTR  = 2
	RoundRectBL  = 4
	RoundRectBR  = 8
	RoundRectALL = RoundRectTL | RoundRectTR | RoundRectBL | RoundRectBR
)

func DrawRoundRect(context *dom.CanvasRenderingContext2D, w, h, r float64, which int) {
	hw := int(w) >> 1
	hh := int(h) >> 1

	if w < 0 || h < 0 {
		return
	}

	if which == 0 {
		which = RoundRectALL
	}

	if (r >= float64(hw) || r >= float64(hh)) && which == RoundRectALL {
		context.Arc(float64(hw), float64(hh), math.Min(float64(hh), float64(hw)), 0, math.Pi*2, false)
		return
	}

	if r > 0 {
		if which&RoundRectTL != 0 {
			context.Arc(r, r, r, math.Pi, 1.5*math.Pi, false)
		} else {
			context.MoveTo(0, 0)
		}

		if which&RoundRectTR != 0 {
			context.LineTo(w-r, 0)
			context.Arc(w-r, r, r, 1.5*math.Pi, 2*math.Pi, false)
		} else {
			context.LineTo(w, 0)
		}

		if which&RoundRectBR != 0 {
			context.LineTo(w, h-r)
			context.Arc(w-r, h-r, r, 0, 0.5*math.Pi, false)
		} else {
			context.LineTo(w, h)
		}

		if which&RoundRectBL != 0 {
			context.LineTo(r, h)
			context.Arc(r, h-r, r, 0.5*math.Pi, math.Pi, false)
		} else {
			context.LineTo(0, h)
		}

		if which&RoundRectTL != 0 {
			context.LineTo(0, r)
		} else {
			context.LineTo(0, 0)
		}
	} else {
		context.Rect(0, 0, w, h)
	}

	return
}
