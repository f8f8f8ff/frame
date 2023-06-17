package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MouseHover struct {
	image.Point
}

type MousePress struct {
	image.Point
	ebiten.MouseButton
}

type MouseRelease struct {
	image.Point
	ebiten.MouseButton
}

type MouseDrag struct {
	Start image.Point
	End   image.Point
	ebiten.MouseButton
	Started  bool
	Released bool
}

func (e *MouseDrag) Update() {
	if e.Released {
		return
	}
	if !MouseJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	if !e.Started {
		e.Started = true
		e.Start = MousePos()
		return
	}
	e.End = MousePos()
	if MouseJustReleased(ebiten.MouseButtonLeft) {
		e.Released = true
	}
}

func (e *MouseDrag) Diff() image.Point {
	return e.End.Sub(e.Start)
}

func (e *MouseDrag) Moved() bool {
	const moveThreshold int = 2
	d := e.Diff()
	if d.X < moveThreshold || d.Y < moveThreshold {
		return false
	}
	return true
}

func (e *MouseDrag) Rect() image.Rectangle {
	return image.Rectangle{
		Min: e.Start,
		Max: e.End,
	}.Canon()
}

func MousePos() image.Point {
	x, y := ebiten.CursorPosition()
	return image.Point{x, y}
}

func MouseJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}
func MouseJustReleased(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}
