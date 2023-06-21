package ui

import (
	"fmt"
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
	case *FlattenOp:
		return &FlattenOp{}
	case *DeleteOp:
		return &DeleteOp{}
	case *LockOrderOp:
		return &LockOrderOp{}
	case *ReorderOp:
		return &ReorderOp{command: op.command}
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
	for _, sp := range c.GetSprites() {
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

type SelectSprOrRectOp struct {
	clr      color.Color
	selDrag  MouseDrag
	target   *sprite.Sprite
	rect     image.Rectangle
	done     bool
	isSprite bool
}

func (op SelectSprOrRectOp) String() string { return "select sprite or rect" }

func (op *SelectSprOrRectOp) Update(ui *UI) (done bool, err error) {
	op.selDrag.Update()
	op.target = ui.Canvas.SpriteAt(MousePos())
	if !op.selDrag.Started {
		return false, nil
	}
	if !op.selDrag.Released {
		return false, nil
	}
	op.done = true
	if op.selDrag.Moved() {
		op.target = nil
		op.isSprite = false
		op.rect = op.selDrag.Rect()
		return true, nil
	}
	op.target = ui.Canvas.SpriteAt(op.selDrag.End)
	op.isSprite = true
	return true, nil
}

func (op *SelectSprOrRectOp) Draw(dst *ebiten.Image) {
	if op.clr == nil {
		op.clr = color.Black
	}
	if !op.selDrag.Started {
		if op.target != nil {
			op.target.Outline(dst, op.clr, 1, -1)
		}
	}
	if op.selDrag.Moved() {
		draw.StrokeRect(dst, op.selDrag.Rect(), op.clr, 1, 0)
	}
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
	sprOrRect *SelectSprOrRectOp
	dstDrag   MouseDrag
	Target    *sprite.Sprite
}

func (op ReshapeOp) String() string { return "reshape" }

func (op *ReshapeOp) Update(ui *UI) (done bool, err error) {
	if op.Target == nil {
		if op.sprOrRect == nil {
			op.sprOrRect = &SelectSprOrRectOp{clr: color.RGBA{0, 0, 255, 255}}
			ui.addOperation(op.sprOrRect)
			return false, nil
		}
		if !op.sprOrRect.done {
			return false, nil
		}
		if op.sprOrRect.isSprite {
			op.Target = op.sprOrRect.target
			if op.Target == nil {
				return true, nil
			}
			return false, nil
		}
		s := ui.Canvas.NewSpriteFromRegion(op.sprOrRect.rect)
		if s == nil {
			return true, nil
		}
		op.Target = s
		ui.Canvas.AddSprite(s)
	}
	if !op.dstDrag.Update() {
		return false, nil
	}
	if !op.dstDrag.Moved() {
		return true, nil
	}
	op.Target.Reshape(op.dstDrag.Rect())
	return true, nil
}

func (op *ReshapeOp) Draw(dst *ebiten.Image) {
	clr := color.RGBA{0, 0, 255, 255}
	if op.Target != nil {
		// outline
		draw.StrokeRect(dst, op.Target.Rect(), clr, 1, -1)
	}
	if !op.dstDrag.Started {
		return
	}
	if op.dstDrag.Moved() {
		opts := draw.ReshapeOpts(op.Target.Rect(), op.dstDrag.Rect())
		op.Target.DrawWithOps(dst, &opts, 1)
	}
	draw.StrokeRect(dst, op.dstDrag.Rect(), clr, 2, -2)
}

type FlattenOp struct {
	drag MouseDrag
	rect image.Rectangle
	spr  *sprite.Sprite
	done bool
}

func (op FlattenOp) String() string { return "liftshape" }

func (op *FlattenOp) Update(ui *UI) (done bool, err error) {
	if op.rect.Empty() {
		if !op.drag.Update() {
			return false, nil
		}
		op.done = true
		if !op.drag.Moved() {
			return true, nil
		}
		op.rect = op.drag.Rect()
	}
	if op.spr == nil {
		op.spr = ui.Canvas.NewSpriteFromRegion(op.rect)
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

type ReorderOp struct {
	selOp   *SingleSelectOp
	target  *sprite.Sprite
	command sprite.ReorderCommand
}

func (op ReorderOp) String() string {
	return fmt.Sprintf("reorder: %v", op.command)
}

func (op *ReorderOp) Update(ui *UI) (done bool, err error) {
	if op.command == sprite.ReorderNone {
		return true, nil
	}
	if op.target == nil {
		if op.selOp == nil {
			op.selOp = &SingleSelectOp{}
			ui.addOperation(op.selOp)
		}
		if !op.selOp.done {
			return false, nil
		}
		if op.selOp.target == nil {
			return true, nil
		}
		op.target = op.selOp.target
	}
	ui.Canvas.Reorder(op.command, op.target)
	return true, nil
}

type SingleSelectOp struct {
	target *sprite.Sprite
	done   bool
	clr    color.Color
}

func (op SingleSelectOp) String() string { return "select(single)" }

func (op *SingleSelectOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.Black
	}
	op.target = ui.Canvas.SpriteAt(MousePos())
	if MouseJustPressed(ebiten.MouseButtonLeft) {
		op.done = true
		return true, nil
	}
	return false, nil
}

func (op *SingleSelectOp) Draw(dst *ebiten.Image) {
	if op.target != nil {
		op.target.Outline(dst, op.clr, 1, -1)
	}
}
