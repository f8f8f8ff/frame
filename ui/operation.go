package ui

import (
	"frame/draw"
	"frame/sprite"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Operation interface {
	Update(*UI) (done bool, err error)
}

type Drawable interface {
	Draw(*ebiten.Image)
}

type SelectOp struct {
	clr     color.Color
	drag    MouseDrag
	Targets []*sprite.Sprite
	moved   bool
	done    bool
}

func (op *SelectOp) Update(ui *UI) (done bool, err error) {
	op.drag.Update()
	if !op.drag.Started {
		return false, nil
	}
	if !op.moved {
		op.moved = op.drag.Moved()
	}
	c := ui.Canvas
	// first click
	if op.drag.JustStarted {
		sp := c.SpriteAt(MousePos())
		if sp == nil {
			return false, nil
		}
		op.Targets = append(op.Targets, sp)
		op.done = true
		return true, nil
	}
	if op.drag.Released {
		op.done = true
		return true, nil
	}
	if !op.moved {
		return false, nil
	}
	op.Targets = []*sprite.Sprite{}
	for _, sp := range c.Sprites() {
		if sp.Overlaps(op.drag.Rect()) {
			op.Targets = append(op.Targets, sp)
		}
	}
	return false, nil
}

func (op *SelectOp) Draw(dst *ebiten.Image) {
	// TODO hover
	for _, sp := range op.Targets {
		// outline
		draw.StrokeRect(dst, sp.Rect(), op.clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 1, 0)
}

type MoveOp struct {
	selOp   *SelectOp
	drag    MouseDrag
	Targets []*sprite.Sprite
}

func (op *MoveOp) Update(ui *UI) (done bool, err error) {
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectOp{clr: color.Black}
			ui.addOperation(op.selOp)
		}
		if !op.selOp.done {
			return false, nil
		}
		op.Targets = op.selOp.Targets
		if len(op.Targets) == 0 {
			return true, nil
		}
	}
	if !op.drag.Update() {
		return false, nil
	}
	// TODO reorder
	//c := ui.Canvas
	// if ui.moveReorder {
	// 	c.RemoveSprite(op.target)
	// 	c.AddSprite(op.target)
	// }
	for _, sp := range op.Targets {
		sp.MoveBy(op.drag.Diff())
	}
	return true, nil
}

func (op *MoveOp) Draw(dst *ebiten.Image) {
	if op.selOp == nil {
		return
	}
	if !op.selOp.done {
		return
	}
	for _, sp := range op.selOp.Targets {
		// TODO outline func
		draw.StrokeRect(dst, sp.Rect(), color.Black, 1, -1)
		if op.drag.Started {
			sp.Draw(dst, op.drag.Diff(), 1)
		}
	}
}

type CropOp struct {
	selOp   *SelectOp
	drag    MouseDrag
	Targets []*sprite.Sprite
}

func (op *CropOp) Update(ui *UI) (done bool, err error) {
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectOp{clr: color.Black}
			ui.addOperation(op.selOp)
		}
		if !op.selOp.done {
			return false, nil
		}
		op.Targets = op.selOp.Targets
		if len(op.Targets) == 0 {
			return true, nil
		}
	}
	if !op.drag.Update() {
		return false, nil
	}
	for _, sp := range op.Targets {
		ui.Canvas.RemoveSprite(sp)
		s := sp.Crop(op.drag.Rect())
		if s == nil {
			continue
		}
		ui.Canvas.AddSprite(s)
	}
	return true, nil
}

func (op *CropOp) Draw(dst *ebiten.Image) {
	// TODO hover
	clr := color.RGBA{0, 255, 0, 255}
	for _, sp := range op.Targets {
		// outline
		draw.StrokeRect(dst, sp.Rect(), clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), clr, 2, 2)
}
