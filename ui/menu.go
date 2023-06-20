package ui

import (
	"frame/draw"
	"frame/sprite"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func MainMenu(ui *UI) *Menu {
	utilityMenuOps := []*MenuOption{
		{text: "(un)lock order", operation: &LockOrderOp{}},
		{text: "delete all", operation: &DeleteOp{Targets: ui.Canvas.Sprites()}},
	}
	utilityMenu := NewMenu(utilityMenuOps, ebiten.MouseButtonLeft)

	options := []*MenuOption{
		{text: "nil", operation: nil},
		{text: "move", operation: &MoveOp{}},
		{text: "crop", operation: &CropOp{}},
		{text: "reshape", operation: &ReshapeOp{}},
		{text: "liftshape", operation: &FlatReshapeOp{}},
		{text: "delete", operation: &DeleteOp{}},
		{text: "util", operation: utilityMenu},
	}
	p := true
	return &Menu{
		options:      options,
		mouseButton:  ebiten.MouseButtonRight,
		startPressed: &p,
	}
}

func NewMenu(opts []*MenuOption, button ebiten.MouseButton) *Menu {
	return &Menu{
		options:     opts,
		mouseButton: button,
	}
}

type MenuOption struct {
	text      string
	operation Operation
	*sprite.Sprite
}

type Menu struct {
	options      []*MenuOption
	mouseButton  ebiten.MouseButton
	rect         *image.Rectangle
	startPressed *bool
	result       Operation
}

func (m *Menu) Update(ui *UI) (done bool, err error) {
	if m.rect == nil {
		m.rect = &image.Rectangle{}
		m.rect.Min = MousePos()
	}
	if m.startPressed == nil {
		pressed := ebiten.IsMouseButtonPressed(m.mouseButton)
		m.startPressed = &pressed
	}
	if m.options[0].Sprite == nil {
		m.createOptionSprites()
		return false, nil
	}
	if *m.startPressed && !MouseJustReleased(m.mouseButton) {
		return false, nil
	} else if !*m.startPressed && !MouseJustPressed(m.mouseButton) {
		return false, nil
	}
	for _, opt := range m.options {
		if opt.In(MousePos()) {
			m.result = opt.operation
			return true, nil
		}
	}
	return true, nil
}

func (m *Menu) createOptionSprites() {
	const height int = 18
	const padding int = 4
	width := 0
	for _, opt := range m.options {
		w := draw.BoundString(draw.Font, opt.text).Dx()
		if w > width {
			width = w
		}
	}
	width += padding * 2
	buttonRect := image.Rect(0, 0, width, height)
	fg := image.Black
	bg := image.White
	for index, opt := range m.options {
		im := draw.NewTextImage(opt.text, draw.Font, buttonRect, padding, fg, bg)
		sp := &sprite.Sprite{
			Image: im,
			Pos:   image.Point{0, index * height},
		}
		sp.Pos = sp.Pos.Add(m.rect.Min)
		opt.Sprite = sp
	}
	overallsize := image.Point{width, height * len(m.options)}
	m.rect.Max = m.rect.Min.Add(overallsize)
}

func (m *Menu) Draw(dst *ebiten.Image) {
	if m.rect == nil {
		return
	}
	for _, opt := range m.options {
		if opt.In(MousePos()) {
			opt.DrawInverted(dst, image.Point{0, 0}, 1)
			continue
		}
		opt.Draw(dst, image.Point{0, 0}, 1)
	}
	// outline menu, invert highlighed
	draw.StrokeRect(dst, *m.rect, color.Black, 2, 2)
}
