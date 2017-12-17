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

func DrawNightPatchEx(context *dom.CanvasRenderingContext2D, image *dom.HTMLImageElement, s_x, s_y, s_w, s_h, x, y, w, h float64) {
	if image == nil {
		context.FillRect(x, y, w, h)
		return
	}

	if s_w == 0 || int(s_w) > image.Width {
		s_w = float64(image.Width)
	}

	if s_h == 0 || int(s_h) > image.Height {
		s_h = float64(image.Height)
	}

	if w < s_w && h < s_h && (s_w < 3 || s_h < 3) {
		context.Call("drawImage", image, s_x, s_y, s_w, s_h, x, y, w, h)
		return
	}

	tw := 0.0
	th := 0.0
	cw := 0.0
	ch := 0.0
	dcw := 0.0
	dch := 0.0

	if w < s_w {
		tw = w / 2
		dcw = 0
		cw = 0
	} else {
		tw = math.Floor(s_w / 3)
		dcw = w - tw - tw
		cw = s_w - tw - tw
	}

	if h < s_h {
		th = h / 2
		dch = 0
		ch = 0
	} else {
		th = math.Floor(s_h / 3)
		dch = h - th - th
		ch = s_h - th - th
	}

	//draw four corner
	context.Call("drawImage", image, s_x, s_y, tw, th, x, y, tw, th)
	context.Call("drawImage", image, s_x+s_w-tw, s_y, tw, th, x+w-tw, y, tw, th)
	context.Call("drawImage", image, s_x, s_y+s_h-th, tw, th, x, y+h-th, tw, th)
	context.Call("drawImage", image, s_x+s_w-tw, s_y+s_h-th, tw, th, x+w-tw, y+h-th, tw, th)

	//top/bottom center
	if dcw > 0 {
		context.Call("drawImage", image, s_x+tw, s_y, cw, th, x+tw, y, dcw, th)
		context.Call("drawImage", image, s_x+tw, s_y+s_h-th, cw, th, x+tw, y+h-th, dcw, th)
	}

	//left/right center
	if dch > 0 {
		context.Call("drawImage", image, s_x, s_y+th, tw, ch, x, y+th, tw, dch)
		context.Call("drawImage", image, s_x+s_w-tw, s_y+th, tw, ch, x+w-tw, y+th, tw, dch)
	}

	if dcw > 0 && dch > 0 {
		context.Call("drawImage", image, s_x+tw, s_y+th, cw, ch, x+tw, y+th, dcw, dch)
	}

	return
}
