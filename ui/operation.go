package ui

import (
	"frame/draw"
	"frame/sprite"
	"image"
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
		op.Targets = []*sprite.Sprite{}
		if sp := ui.Canvas.SpriteAt(MousePos()); sp != nil {
			op.Targets = append(op.Targets, sp)
		}
		return false, nil
	}
	if !op.moved {
		op.moved = op.drag.Moved()
	}
	c := ui.Canvas
	// first click
	if op.drag.JustStarted {
		op.Targets = []*sprite.Sprite{}
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
	for _, sp := range op.Targets {
		if !ui.LockOrder {
			ui.Canvas.RemoveSprite(sp)
			ui.Canvas.AddSprite(sp)
		}
		sp.MoveBy(op.drag.Diff())
	}
	return true, nil
}

func (op *MoveOp) Draw(dst *ebiten.Image) {
	if len(op.Targets) == 0 {
		return
	}
	for _, sp := range op.Targets {
		// TODO outline func
		draw.StrokeRect(dst, sp.Rect(), color.Black, 1, -1)
		if op.drag.Started {
			sp.Draw(dst, op.drag.Diff(), 1)
		}
	}
}

type DragOp struct {
	moveOp  *MoveOp
	Targets []*sprite.Sprite
}

func (op *DragOp) Update(ui *UI) (done bool, err error) {
	if len(op.Targets) == 0 {
		if MouseJustPressed(ebiten.MouseButtonLeft) {
			s := ui.Canvas.SpriteAt(MousePos())
			if s == nil {
				return true, nil
			}
			op.Targets = append(op.Targets, s)
		}
		// return false, nil
	}
	op.moveOp = &MoveOp{
		Targets: op.Targets,
	}
	// update once to update the drag on the same update
	if done, _ := op.moveOp.Update(ui); !done {
		ui.addOperation(op.moveOp)
	}
	return true, nil
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

type ReshapeOp struct {
	drag   MouseDrag
	Target *sprite.Sprite
}

func (op *ReshapeOp) Update(ui *UI) (done bool, err error) {
	if op.Target == nil {
		if MouseJustPressed(ebiten.MouseButtonLeft) {
			s := ui.Canvas.SpriteAt(MousePos())
			if s == nil {
				return true, nil
			}
			op.Target = s
		}
		return false, nil
	}
	if !op.drag.Update() {
		return false, nil
	}
	if !op.drag.Moved() {
		return true, nil
	}
	op.Target.Reshape(op.drag.Rect())
	return true, nil
}

func (op *ReshapeOp) Draw(dst *ebiten.Image) {
	clr := color.RGBA{0, 0, 255, 255}
	if op.Target != nil {
		// outline
		draw.StrokeRect(dst, op.Target.Rect(), clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	if op.drag.Moved() {
		opts := draw.ReshapeOpts(op.Target.Rect(), op.drag.Rect())
		op.Target.DrawWithOps(dst, &opts, 1)
	}
	draw.StrokeRect(dst, op.drag.Rect(), clr, 2, -2)
}

type FlatReshapeOp struct {
	drag  MouseDrag
	drag2 MouseDrag
	spr   *sprite.Sprite
}

func (op *FlatReshapeOp) Update(ui *UI) (done bool, err error) {
	if !op.drag.Update() {
		return false, nil
	}
	if !op.drag.Moved() {
		return true, nil
	}

	if op.spr == nil {
		im, r := draw.CropImage(ui.Canvas.Image(), op.drag.Rect(), image.Point{0, 0})
		if im == nil {
			return true, nil // error
		}
		op.spr = &sprite.Sprite{
			Image: im,
			Pos:   r.Min,
		}
	}

	if !op.drag2.Update() {
		return false, nil
	}
	if !op.drag2.Moved() {
		return true, nil
	}

	op.spr.Reshape(op.drag2.Rect())
	ui.Canvas.AddSprite(op.spr)
	return true, nil
}

func (op *FlatReshapeOp) Draw(dst *ebiten.Image) {
	clr := color.RGBA{0, 0, 255, 255}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), clr, 1, 1)
	if !op.drag2.Started {
		return
	}
	if op.drag2.Moved() {
		opts := draw.ReshapeOpts(op.spr.Rect(), op.drag2.Rect())
		op.spr.DrawWithOps(dst, &opts, 1)
	}
	draw.StrokeRect(dst, op.drag2.Rect(), clr, 2, -2)
}

type DeleteOp struct {
	selOp   *SelectOp
	Targets []*sprite.Sprite
}

func (op *DeleteOp) Update(ui *UI) (done bool, err error) {
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectOp{clr: color.RGBA{255, 0, 0, 255}}
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
	for _, sp := range op.Targets {
		ui.Canvas.RemoveSprite(sp)
	}
	return true, nil
}

type LockOrderOp struct{}

func (op *LockOrderOp) Update(ui *UI) (done bool, err error) {
	ui.LockOrder = !ui.LockOrder
	return true, nil
}
