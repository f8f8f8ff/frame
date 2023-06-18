package ui

import (
	"frame/draw"
	"frame/sprite"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Menu struct {
	options []*struct {
		text      string
		operation interface{}
		*sprite.Sprite
	}
	mouseButton ebiten.MouseButton
	rect        image.Rectangle
}

func MainMenu() *Menu {
	options := []*struct {
		text      string
		operation interface{}
		*sprite.Sprite
	}{
		{text: "test", operation: nil},
		{text: "2", operation: nil},
		{text: "select", operation: &SelectOp{clr: color.Black}},
	}
	r := image.Rectangle{}
	r.Min = MousePos()
	return &Menu{
		options:     options,
		mouseButton: ebiten.MouseButtonRight,
		rect:        r,
	}
}

func (m *Menu) Update() (op interface{}, done bool) {
	if m.options[0].Sprite == nil {
		m.createOptionSprites()
		return nil, false
	}
	if !MouseJustReleased(m.mouseButton) {
		return nil, false
	}
	for _, opt := range m.options {
		if opt.In(MousePos()) {
			return opt.operation, true
		}
	}
	return nil, true
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
	for _, opt := range m.options {
		if opt.In(MousePos()) {
			opt.DrawInverted(dst, image.Point{0, 0}, 1)
			continue
		}
		opt.Draw(dst, image.Point{0, 0}, 1)
	}
	// outline menu, invert highlighed
	draw.StrokeRect(dst, m.rect, color.Black, 2, 2)
}
