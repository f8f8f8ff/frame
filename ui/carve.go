package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"frame/draw"
	"frame/sprite"

	"carve"
)

type CarveOp struct {
	sprOrRect *SelectSpriteRectOp
	dstDrag   MouseDrag
	Target    *sprite.Sprite
	clr       color.Color
}

func (op CarveOp) String() string { return "carve" }

func (op *CarveOp) Update(ui *UI) (done bool, err error) {
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
	// op.Target.Reshape(op.dstDrag.Rect())
	img := carve.Resize(op.Target.Image, op.dstDrag.Rect().Dx(), op.dstDrag.Rect().Dy())
	s := &sprite.Sprite{
		Image:         ebiten.NewImageFromImage(img),
		Pos:           op.dstDrag.Rect().Min,
		OpacityOffset: 0,
	}
	ui.Canvas.RemoveSprite(op.Target)
	ui.Canvas.AddSprite(s)
	return true, nil
}

func (op *CarveOp) Draw(dst *ebiten.Image) {
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
