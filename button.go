package gwk

type Button struct {
	*ImageText
}

func NewButton(parent *Widget, x, y, w, h float32) *Button {
	button := &Button{
		ImageText: NewImageText(parent, x, y, w, h),
	}
	button.t = TYPE_BUTTON
	button.I = button

	return button
}

func NewOkButton(parent *Widget, x, y, w, h float32) *Button {
	button := NewButton(parent, x, y, w, h)
	button.UseTheme("button-ok")

	return button
}

func NewCancelButton(parent *Widget, x, y, w, h float32) *Button {
	button := NewButton(parent, x, y, w, h)
	button.UseTheme("button-cancel")

	return button
}
