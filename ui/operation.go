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
}

func (op *SelectOp) Update(ui *UI) (done bool, err error) {
	if CancelInput() {
		return true, nil
	}
	op.drag.Update()
	if !op.drag.Started {
		return false, nil
	}
	if !op.moved {
		op.moved = op.drag.Moved()
	}
	c := ui.Canvas
	if len(op.Targets) == 0 && !op.moved {
		sp := c.SpriteAt(MousePos())
		if sp == nil {
			return false, nil
		}
		op.Targets = append(op.Targets, sp)
		return true, nil
	}
	op.Targets = []*sprite.Sprite{}
	for _, sp := range c.Sprites() {
		if sp.Overlaps(op.drag.Rect()) {
			op.Targets = append(op.Targets, sp)
		}
	}
	if !op.drag.Released {
		return false, nil
	}
	return true, nil
}

func (op *SelectOp) Draw(dst *ebiten.Image) {
	for _, sp := range op.Targets {
		draw.StrokeRect(dst, sp.Rect(), op.clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 1, 0)
}

type DragOp struct {
	drag   MouseDrag
	target *sprite.Sprite
}

func (op *DragOp) Update(ui *UI) (done bool, err error) {
	op.drag.Update()
	if !op.drag.Started {
		return false, nil
	}
	c := ui.Canvas
	if op.target == nil {
		op.target = c.SpriteAt(MousePos())
		if op.target == nil {
			return true, nil
		}
	}
	if ui.moveReorder {
		c.RemoveSprite(op.target)
		c.AddSprite(op.target)
	}
	if !op.drag.Released {
		return false, nil
	}
	op.target.MoveBy(op.drag.Diff())
	return true, nil
}
