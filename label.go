package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/theme"
	"honnef.co/go/js/dom"
	"math"
)

type Label struct {
	*Widget
	textU        bool
	textI        bool
	textB        bool
	font         string
	leftBorder   int
	rightBorder  int
	topBorder    int
	bottomBorder int
	fontSize     int
	textColor    string
	lineColor    string
	textAlignH   string
	textAlignV   string
	flexibleSize int
	singleLine   bool
	lines        []string
}

func NewLabel(parent *Widget, x, y, w, h float32) *Label {
	label := &Label{
		Widget:       NewWidget(TYPE_LABEL, parent, x, y, w, h),
		flexibleSize: 4,
		fontSize:     14,
		leftBorder:   2,
		rightBorder:  2,
		topBorder:    2,
		bottomBorder: 2,
		singleLine:   true,
		textAlignV:   "middle",
		textAlignH:   "center",
	}
	label.I = label

	return label
}

func (label *Label) getTipsStyle() *theme.ThemeStyle {
	return label.theme.StateNormal
}

func (label *Label) drawTips(context *dom.CanvasRenderingContext2D) *Label {
	if label.state != STATE_OVER {
		return label
	}

	tips := label.getTips()
	if len(tips) > 0 {
		style := label.getTipsStyle()
		h := 30
		x := label.getWidth() + 3
		y := label.getHeight()
		w := context.MeasureText(tips).Width + 40

		context.LineWidth = 1
		context.FillStyle = style.TipsFillColor
		context.StrokeStyle = style.TipsLineColor

		context.Rect(float64(x), float64(y), w, float64(h))
		context.Fill()
		context.Stroke()

		context.TextAlign = "center"
		context.TextBaseline = "middle"
		context.Font = style.Font
		context.FillStyle = style.TipsTextColor
		x = x + (int(w) >> 1)
		y = y + (h >> 1)
		context.FillText(tips, float64(x), float64(y), -1)
	}

	return label
}

func (label *Label) setBorder(sides ...int) *Label {
	if len(sides) > 0 {
		label.leftBorder = sides[0]
	}

	if len(sides) > 1 {
		label.topBorder = sides[1]
	}

	if len(sides) > 2 {
		label.rightBorder = sides[2]
	}

	if len(sides) > 3 {
		label.bottomBorder = sides[3]
	}

	return label
}

func (label *Label) setLayoutFlexibleSize(flexibleSize int) *Label {
	label.flexibleSize = flexibleSize

	return label
}

func (label *Label) layoutText(context *dom.CanvasRenderingContext2D, text string) {
	width := label.rect.W - label.leftBorder - label.rightBorder
	if len(text) > 0 {
		context.Font = label.getFont()
		label.lines = layoutText(context, label.fontSize, text, width, label.flexibleSize)
	} else {
		label = nil
	}

	return
}

func (label *Label) relayout(context *dom.CanvasRenderingContext2D, force bool) {
	if !label.needRelayout && !force && context == nil {
		return
	}

	text := label.GetText()
	label.layoutText(context, text)
	label.needRelayout = false

	return
}

func (label *Label) SetText(str string, notify bool) *Label {
	label.text = str

	if notify && label.onChanged != nil {
		label.onChanged(label.text)
	}
	label.setNeedRelayout(true)

	return label
}

func (label *Label) setTextAlignV(align string) *Label {
	label.textAlignV = align

	return label
}

func (label *Label) setTextAlignH(align string) *Label {
	label.textAlignH = align

	return label
}

func (label *Label) setTextColor(textColor string) *Label {
	label.textColor = textColor

	return label
}

func (label *Label) setLineColor(lineColor string) *Label {
	label.lineColor = lineColor

	return label
}

func (label *Label) getTextColor() string {
	if len(label.textColor) != 0 {
		return label.textColor
	} else {
		return label.getStyle("").TextColor
	}
}

func (label *Label) getLineColor() string {
	if len(label.lineColor) != 0 {
		return label.lineColor
	} else {
		return label.getStyle("").LineColor
	}
}

func (label *Label) setTextBold(textB bool) *Label {
	label.textB = textB

	return label.updateFont()
}

func (label *Label) setTextUnderline(textUnderline bool) *Label {
	label.textU = textUnderline

	return label.updateFont()
}

func (label *Label) setTextItalic(textItalic bool) *Label {
	label.textI = textItalic

	return label.updateFont()
}

func (label *Label) setFontSize(fontSize int) *Label {
	label.fontSize = fontSize

	return label.updateFont()
}

func (label *Label) setSingleLineMode(singleLine bool) *Label {
	label.singleLine = singleLine

	return label
}

func (label *Label) updateFont() *Label {
	label.font = fmt.Sprint(label.fontSize) + "px"

	if label.textB {
		label.font += "bold "
	}

	if label.textI {
		label.font += "italic "
	}

	label.font += "sans-serif"

	return label
}

func (label *Label) getFont() string {
	if len(label.font) > 0 {
		return label.font
	} else {
		style := label.getStyle("")
		return style.Font
	}
}

func (label *Label) getLines() []string {
	return label.lines
}

func (label *Label) paintSelf(context *dom.CanvasRenderingContext2D) {
	if label.singleLine {
		label.paintSelfSL(context)
	} else {
		label.paintSelfML(context)
	}

	return
}

func (label *Label) paintSelfSL(context *dom.CanvasRenderingContext2D) {
	text := label.text
	label.paintSLText(context, text)

	return
}

func (label *Label) paintSLText(context *dom.CanvasRenderingContext2D, text string) {
	context.Font = label.getFont()
	context.TextBaseline = "middle"
	context.FillStyle = label.getTextColor()

	var x int
	var y = label.getHeight() >> 1
	var w = label.getWidth()

	switch label.textAlignH {
	case "center":
		x = w >> 1
		context.TextAlign = "center"
	case "right":
		x = w - label.rightBorder
		context.TextAlign = "right"
	default:
		x = label.leftBorder
		context.TextAlign = "left"
	}
	context.FillText(text, float64(x), float64(y), float64(w))

	return
}

func (label *Label) paintSelfML(context *dom.CanvasRenderingContext2D) {
	lines := label.getLines()
	if len(lines) == 0 {
		return
	}

	fontSize := label.fontSize
	lineHeight := float64(fontSize) * 1.5
	textHeight := lineHeight * float64(len(lines))
	height := label.rect.H - label.topBorder - label.bottomBorder
	maxLineNr := int(math.Min(float64(len(lines)), math.Floor(float64(height)/lineHeight)))

	var x, y int

	textHeight = float64(maxLineNr) * lineHeight
	switch label.textAlignV {
	case "middle":
		y = (label.rect.H - int(textHeight)) >> 1
	case "bottom":
		y = label.rect.H - int(textHeight) - label.bottomBorder
	default:
		y = label.bottomBorder
	}

	width := label.rect.W
	leftBorder := label.leftBorder
	rightBorder := label.rightBorder

	context.TextAlign = "left"
	context.TextBaseline = "top"
	context.StrokeStyle = label.getLineColor()
	context.FillStyle = label.getTextColor()
	context.Font = label.getFont()
	context.LineWidth = 1

	for i := 0; i < maxLineNr; i++ {
		str := lines[i]
		if len(str) == 0 {
			y += int(lineHeight)
			continue
		}
		textWidth := context.MeasureText(str).Width
		switch label.textAlignH {
		case "center":
			x = (width - int(textWidth)) >> 1
		case "right":
			x = width - rightBorder - int(textWidth)
		default:
			x = leftBorder
		}

		context.BeginPath()

		if label.textU {
			ly := y + fontSize + 4
			context.MoveTo(float64(x), float64(ly))
			context.LineTo(float64(x)+textWidth, float64(ly))
			context.Stroke()
		}

		context.FillText(str, float64(x), float64(y), -1)
		y += int(lineHeight)
	}

	return
}
