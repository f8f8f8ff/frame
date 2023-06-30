package ui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"frame/canvas"
	"frame/draw"
	"frame/sprite"
)

type Operation interface {
	Update(*UI) (done bool, err error)
}

type Drawable interface {
	Draw(*ebiten.Image)
}

type FullDrawer interface {
	FullDraw(*ebiten.Image, *canvas.Canvas)
}

// definies how to repeat an operation
func CopyOperation(op Operation) Operation {
	switch op := op.(type) {
	case *SelectSpriteMultiOp:
		return &SelectSpriteMultiOp{clr: op.clr}
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
	case *CopyOp:
		return &CopyOp{}
	case *CutOp:
		return &CutOp{}
	case *CarveOp:
		return &CarveOp{}
	}
	return nil
}

type SelectSpriteMultiOp struct {
	clr     color.Color
	drag    MouseDrag
	Targets []*sprite.Sprite
	moved   bool
	done    bool
}

func (op SelectSpriteMultiOp) String() string { return "select sprite(s)" }

func (op *SelectSpriteMultiOp) Update(ui *UI) (done bool, err error) {
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

func (op *SelectSpriteMultiOp) Draw(dst *ebiten.Image) {
	if op.clr == nil {
		op.clr = color.Black
	}
	for _, sp := range op.Targets {
		sp.Outline(dst, op.clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 1, 0)
}

type SelectSpriteRectOp struct {
	clr      color.Color
	selDrag  MouseDrag
	target   *sprite.Sprite
	rect     image.Rectangle
	done     bool
	isSprite bool
}

func (op SelectSpriteRectOp) String() string { return "select sprite or region" }

func (op *SelectSpriteRectOp) Update(ui *UI) (done bool, err error) {
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

func (op *SelectSpriteRectOp) Draw(dst *ebiten.Image) {
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

type SelectSpriteOp struct {
	target *sprite.Sprite
	done   bool
	clr    color.Color
}

func (op SelectSpriteOp) String() string { return "select sprite" }

func (op *SelectSpriteOp) Update(ui *UI) (done bool, err error) {
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

func (op *SelectSpriteOp) Draw(dst *ebiten.Image) {
	if op.clr == nil {
		op.clr = color.Black
	}
	if op.target != nil {
		op.target.Outline(dst, op.clr, 1, -1)
	}
}

type MoveOp struct {
	selOp   *SelectSpriteMultiOp
	drag    MouseDrag
	Targets []*sprite.Sprite
	clr     color.Color
}

func (op MoveOp) String() string { return "move" }

func (op *MoveOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.Black
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
		if op.drag.Started {
			sp.Draw(dst, op.drag.Diff(), 1)
		}
		sp.Outline(dst, op.clr, 1, -1)
	}
}

type DragOp struct {
	moveOp  *MoveOp
	Targets []*sprite.Sprite
}

func (op DragOp) String() string { return "drag" }

func (op *DragOp) Update(ui *UI) (done bool, err error) {
	// TODO update to use SelectSpriteOp
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
	selOp   *SelectSpriteMultiOp
	drag    MouseDrag
	Targets []*sprite.Sprite
	clr     color.Color
}

func (op CropOp) String() string { return "crop" }

func (op *CropOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{255, 255, 0, 255} // yellow
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
		i := ui.Canvas.Sprites.IndexOf(sp)
		s := sp.Crop(op.drag.Rect())
		if s == nil {
			ui.Canvas.RemoveSprite(sp)
			continue
		}
		ui.Canvas.Sprites[i] = s
	}
	return true, nil
}

func (op *CropOp) Draw(dst *ebiten.Image) {
	for _, sp := range op.Targets {
		sp.Outline(dst, op.clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 2, 2)
}

type ReshapeOp struct {
	sprOrRect *SelectSpriteRectOp
	dstDrag   MouseDrag
	Target    *sprite.Sprite
	clr       color.Color
}

func (op ReshapeOp) String() string { return "reshape" }

func (op *ReshapeOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{0, 0, 255, 255}
	}
	if op.Target == nil {
		if op.sprOrRect == nil {
			op.sprOrRect = &SelectSpriteRectOp{clr: op.clr}
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
	if op.Target != nil {
		op.Target.Outline(dst, op.clr, 1, -1)
	}
	if !op.dstDrag.Started {
		return
	}
	if op.dstDrag.Moved() {
		opts := draw.ReshapeOpts(op.Target.Rect(), op.dstDrag.Rect())
		op.Target.DrawWithOps(dst, &opts, 1)
	}
	draw.StrokeRect(dst, op.dstDrag.Rect(), op.clr, 2, -2)
}

type FlattenOp struct {
	drag MouseDrag
	rect image.Rectangle
	spr  *sprite.Sprite
	done bool
	clr  color.Color
}

func (op FlattenOp) String() string { return "flatten" }

func (op *FlattenOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{50, 205, 50, 255} // lime green
	}
	if op.rect.Empty() {
		if !op.drag.Update() {
			return false, nil
		}
		op.done = true
		if !op.drag.Moved() {
			op.rect = ui.Canvas.Image().Bounds()
		} else {
			op.rect = op.drag.Rect()
		}
	}
	if op.spr == nil {
		op.spr = ui.Canvas.NewSpriteFromRegion(op.rect)
	}
	ui.Canvas.AddSprite(op.spr)
	return true, nil
}

func (op *FlattenOp) Draw(dst *ebiten.Image) {
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 1, 1)
}

type DeleteOp struct {
	selOp   *SelectSpriteMultiOp
	Targets []*sprite.Sprite
	clr     color.Color
}

func (op DeleteOp) String() string { return "delete" }

func (op *DeleteOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{255, 0, 0, 255} // red
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
		ui.RemoveSprite(sp)
	}
	return true, nil
}

type DeleteAllOp struct{}

func (op DeleteAllOp) String() string { return "delete all" }

func (op DeleteAllOp) Update(ui *UI) (done bool, err error) {
	ui.Canvas.ClearSprites()
	return true, nil
}

type LockOrderOp struct{}

func (op LockOrderOp) String() string { return "(un)lock order" }

func (op *LockOrderOp) Update(ui *UI) (done bool, err error) {
	ui.LockOrder = !ui.LockOrder
	return true, nil
}

type ReorderOp struct {
	selOp   *SelectSpriteOp
	target  *sprite.Sprite
	command sprite.ReorderCommand
	clr     color.Color
}

func (op ReorderOp) String() string {
	return fmt.Sprintf("reorder: %v", op.command)
}

func (op *ReorderOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.Black
	}
	if op.command == sprite.ReorderNone {
		return true, nil
	}
	if op.target == nil {
		if op.selOp == nil {
			op.selOp = &SelectSpriteOp{clr: op.clr}
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

type CopyOp struct {
	selOp   *SelectSpriteMultiOp
	drag    MouseDrag
	Targets []*sprite.Sprite
	clr     color.Color
}

func (op CopyOp) String() string { return "move" }

func (op *CopyOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{0, 200, 0, 255} // green
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
		copy := sp.Copy()
		copy.MoveBy(op.drag.Diff())
		ui.AddSprite(copy)
	}
	return true, nil
}

func (op *CopyOp) Draw(dst *ebiten.Image) {
	if len(op.Targets) == 0 {
		return
	}
	for _, sp := range op.Targets {
		if op.drag.Started {
			sp.Draw(dst, op.drag.Diff(), 1)
		}
		sp.Outline(dst, op.clr, 1, -1)
	}
}

type CutOp struct {
	selOp   *SelectSpriteMultiOp
	drag    MouseDrag
	Targets []*sprite.Sprite
	clr     color.Color
}

func (op CutOp) String() string { return "cut" }

func (op *CutOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{255, 165, 0, 255} // orange
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
		draw.CutImage(sp.Image, op.drag.Rect().Sub(sp.Pos))
	}
	return true, nil
}

func (op *CutOp) Draw(dst *ebiten.Image) {
	for _, sp := range op.Targets {
		sp.Outline(dst, op.clr, 1, -1)
	}
	if !op.drag.Started {
		return
	}
	draw.StrokeRect(dst, op.drag.Rect(), op.clr, 2, 2)
}

type OpacityOp struct {
	selOp         *SelectSpriteMultiOp
	drag          MouseDrag
	Targets       []*sprite.Sprite
	clr           color.Color
	opacityOffset float64
}

func (op OpacityOp) String() string { return "change opacity" }

func (op *OpacityOp) Update(ui *UI) (done bool, err error) {
	if op.clr == nil {
		op.clr = color.RGBA{0, 0, 0, 255} // black
	}
	if len(op.Targets) == 0 {
		if op.selOp == nil {
			op.selOp = &SelectSpriteMultiOp{clr: op.clr}
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
	released := op.drag.Update()
	if !op.drag.Started {
		return false, nil
	}
	op.opacityOffset = float64(-op.drag.Diff().Y) / 100
	if !released {
		return false, nil
	}

	for _, sp := range op.Targets {
		o := sp.OpacityOffset + op.opacityOffset
		if o > 0 {
			o = 0
		} else if o < -1 {
			o = -1
		}
		sp.OpacityOffset = o
	}
	return true, nil
}

func (op *OpacityOp) FullDraw(dst *ebiten.Image, c *canvas.Canvas) {
	dst.Fill(color.White)
	for i := len(c.Sprites) - 1; i >= 0; i-- {
		sp := c.Sprites[i]
		cont := false
		for _, t := range op.Targets {
			if sp == t {
				cont = true
				break
			}
		}
		if cont {
			continue
		}
		sp.Draw(dst, image.Point{}, 1)
	}
	for _, sp := range op.Targets {
		sp.Outline(dst, op.clr, 1, -1)
	}
	for i := len(op.Targets) - 1; i >= 0; i-- {
		sp := op.Targets[i]
		sp.Draw(dst, image.Point{0, 0}, 1+op.opacityOffset)
	}
}

type CBCopyOp struct {
	flattenOp *FlattenOp
}

func (op *CBCopyOp) String() string { return "copy region to clipboard" }

func (op *CBCopyOp) Update(ui *UI) (done bool, err error) {
	if op.flattenOp == nil {
		op.flattenOp = &FlattenOp{}
		ui.addOperation(op.flattenOp)
		return false, nil
	}
	if op.flattenOp.spr == nil {
		return false, nil
	}
	ui.RemoveSprite(op.flattenOp.spr)
	return true, copyClipboard(op.flattenOp.spr.Image)
}

type CBPasteOp struct {
	spr    *sprite.Sprite
	setPos bool
}

func (op *CBPasteOp) String() string { return "paste from clipboard" }

func (op *CBPasteOp) Update(ui *UI) (done bool, err error) {
	if op.spr == nil {
		op.spr, err = handlePaste()
		if err != nil || op.spr == nil {
			return true, err
		}
		ui.AddSprite(op.spr)
	}
	op.spr.Pos = MousePos()
	if op.setPos || MouseJustPressed(ebiten.MouseButtonLeft) {
		return true, nil
	}
	return false, nil
}

func (op *CBPasteOp) Draw(dst *ebiten.Image) {
	if op.spr == nil {
		return
	}
	op.spr.Draw(dst, image.Point{}, 1)
}
