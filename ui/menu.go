package ui

import (
	"fmt"
	"frame/draw"
	"frame/sprite"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	menuItemHeight = 18
	menuPadding    = 4
	menuFg         = color.Black
	menuBg         = color.White
	menuPaddingClr = color.Black
)

func MainMenu(ui *UI) *Menu {
	reorderMenuOps := []*MenuOption{
		{text: "bring to front", operation: &ReorderOp{command: sprite.ReorderBringToFront}},
		{text: "send to back", operation: &ReorderOp{command: sprite.ReorderSendToBack}},
		{text: "bring forwards", operation: &ReorderOp{command: sprite.ReorderBringForwards}},
		{text: "send backwards", operation: &ReorderOp{command: sprite.ReorderSendBackwards}},
	}
	reorderMenu := NewMenu(reorderMenuOps, ebiten.MouseButtonLeft)
	utilityMenuOps := []*MenuOption{
		{text: "copy to clipboard", operation: &CBCopyOp{}},
		{text: "paste from clipboard", operation: &CBPasteOp{}},
		{text: "(un)lock order", operation: &LockOrderOp{}},
		{text: "delete all", operation: &DeleteAllOp{}},
	}
	utilityMenu := NewMenu(utilityMenuOps, ebiten.MouseButtonLeft)
	options := []*MenuOption{
		RepeatMenuOption(ui.lastOp),
		{text: "move", operation: &MoveOp{}},
		{text: "copy", operation: &CopyOp{}},
		{text: "crop", operation: &CropOp{}},
		{text: "cut", operation: &CutOp{}},
		{text: "reshape", operation: &ReshapeOp{}},
		{text: "carve", operation: &CarveOp{}},
		{text: "flatten", operation: &FlattenOp{}},
		{text: "opacity", operation: &OpacityOp{}},
		{text: "delete", operation: &DeleteOp{}},
		{text: "reorder", operation: reorderMenu},
		{text: "util", operation: utilityMenu},
	}
	p := true
	return &Menu{
		options:      options,
		mouseButton:  ebiten.MouseButtonRight,
		startPressed: &p,
	}
}

func RepeatMenuOption(op Operation) *MenuOption {
	name := fmt.Sprintf("(%v)", op)
	return &MenuOption{
		text:      name,
		operation: op,
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
	screensize   image.Point
}

func (m *Menu) Update(ui *UI) (done bool, err error) {
	if m.rect == nil {
		m.rect = &image.Rectangle{}
		m.rect.Min = MousePos()
	}
	if m.screensize.Eq(image.Pt(0, 0)) {
		m.screensize = image.Pt(ui.Width, ui.Height)
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
	width := 0
	for _, opt := range m.options {
		w := draw.BoundString(draw.Font, opt.text).Dx()
		if w > width {
			width = w
		}
	}
	width += menuPadding * 2
	buttonRect := image.Rect(0, 0, width, menuItemHeight)

	overallsize := image.Point{width, menuItemHeight * len(m.options)}
	m.rect.Max = m.rect.Min.Add(overallsize)
	// +2 for menu border
	if dx := m.rect.Max.X - m.screensize.X + 2; dx > 0 {
		r := m.rect.Add(image.Point{-dx, 0})
		m.rect = &r
	}
	if dy := m.rect.Max.Y - m.screensize.Y + 2; dy > 0 {
		r := m.rect.Add(image.Point{0, -dy})
		m.rect = &r
	}

	for index, opt := range m.options {
		im := draw.NewTextImage(opt.text, draw.Font, buttonRect, menuPadding, menuFg, menuBg)
		sp := &sprite.Sprite{
			Image: im,
			Pos:   image.Point{0, index * menuItemHeight},
		}
		sp.Pos = sp.Pos.Add(m.rect.Min)
		opt.Sprite = sp
	}
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
	draw.StrokeRect(dst, *m.rect, menuPaddingClr, 2, 2)
}
