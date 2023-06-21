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

// definies how to repeat an operation
func CopyOp(op Operation) Operation {
	switch op := op.(type) {
	case *SelectOp:
		return &SelectOp{clr: op.clr}
	case *MoveOp:
		return &MoveOp{}
	case *DragOp:
		return &DragOp{}
	case *CropOp:
		return &CropOp{}
	case *ReshapeOp:
		return &ReshapeOp{}
	case *FlatReshapeOp:
		return &FlatReshapeOp{}
	case *DeleteOp:
		return &DeleteOp{}
	case *LockOrderOp:
		return &LockOrderOp{}
	}
	return nil
}

type SelectOp struct {
	clr     color.Color
	drag    MouseDrag
	Targets []*sprite.Sprite
	moved   bool
	done    bool
}

func (op SelectOp) String() string { return "select" }

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

func (op MoveOp) String() string { return "move" }

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

func (op DragOp) String() string { return "drag" }

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

func (op CropOp) String() string { return "crop" }

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

func (op ReshapeOp) String() string { return "reshape" }

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
	flatOp *FlattenOp
	drag   MouseDrag
	spr    *sprite.Sprite
}

func (op FlatReshapeOp) String() string { return "liftshape" }

func (op *FlatReshapeOp) Update(ui *UI) (done bool, err error) {
	if op.flatOp == nil {
		op.flatOp = &FlattenOp{}
		ui.addOperation(op.flatOp)
		return false, nil
	}
	if !op.flatOp.done {
		return false, nil
	}
	if op.flatOp.spr == nil {
		return true, nil
	}
	if op.spr == nil {
		op.spr = op.flatOp.spr
	}
	if !op.drag.Update() {
		return false, nil
	}
	if !op.drag.Moved() {
		return true, nil
	}
	op.spr.Reshape(op.drag.Rect())
	ui.Canvas.AddSprite(op.spr)
	return true, nil
}

func (op *FlatReshapeOp) Draw(dst *ebiten.Image) {
	clr := color.RGBA{0, 0, 255, 255}
	if op.spr == nil {
		return
	}
	// outline
	draw.StrokeRect(dst, op.spr.Rect(), clr, 1, -1)
	if !op.drag.Started {
		return
	}
	if op.drag.Moved() {
		opts := draw.ReshapeOpts(op.spr.Rect(), op.drag.Rect())
		op.spr.DrawWithOps(dst, &opts, 1)
	}
	draw.StrokeRect(dst, op.drag.Rect(), clr, 2, -2)
}

type FlattenOp struct {
	drag MouseDrag
	spr  *sprite.Sprite
	done bool
}

func (op FlattenOp) String() string { return "liftshape" }

func (op *FlattenOp) Update(ui *UI) (done bool, err error) {
	if !op.drag.Update() {
		return false, nil
	}
	op.done = true
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
	ui.Canvas.AddSprite(op.spr)
	return true, nil
}

func (op *FlattenOp) Draw(dst *ebiten.Image) {
	clr := color.RGBA{0, 0, 255, 255}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), clr, 1, 1)
}

type DeleteOp struct {
	selOp   *SelectOp
	Targets []*sprite.Sprite
}

func (op DeleteOp) String() string { return "delete" }

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

func (op LockOrderOp) String() string { return "(un)lock order" }

func (op *LockOrderOp) Update(ui *UI) (done bool, err error) {
	ui.LockOrder = !ui.LockOrder
	return true, nil
}
