package ui

import (
	"fmt"
	"image"
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
	if op.Target.Rect().Dx() < 5 || op.Target.Rect().Dy() < 5 {
		return true, nil
	}
	if !op.dstDrag.Update() {
		return false, nil
	}
	if !op.dstDrag.Moved() {
		return true, nil
	}

	resultc, progressc := carve.ResizeConc(op.Target.Image, op.dstDrag.Rect().Dx(), op.dstDrag.Rect().Dy())
	go func() {
		var img image.Image
		for img == nil {
			img = <-resultc
		}
		s := &sprite.Sprite{
			Image:         ebiten.NewImageFromImage(img),
			Pos:           op.dstDrag.Rect().Min,
			OpacityOffset: 0,
		}
		ui.Canvas.AddSprite(s)
	}()
	cp := &CarveProgress{
		progress:     0,
		progressChan: progressc,
		dragRect:     op.dstDrag.Rect(),
		clr:          op.clr,
	}
	ui.addOperation(cp)
	return true, nil
}

// TODO make carveop draw rectangle while loading
// will take some work in ui, maybe a new queue of async
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

type CarveProgress struct {
	progress     float64
	progressChan <-chan float64
	dragRect     image.Rectangle
	clr          color.Color
}

func (op *CarveProgress) String() string {
	return fmt.Sprintf("carving: %.0f%%", op.progress*100)
}

func (op *CarveProgress) Update(ui *UI) (done bool, err error) {
	select {
	case p, more := <-op.progressChan:
		op.progress = p
		if !more {
			break
		}
		return false, nil
	default:
		return false, nil
	}
	return true, nil
}

func (op *CarveProgress) Draw(dst *ebiten.Image) {
	dr := op.dragRect
	x1 := dr.Min.X + int(float64(dr.Dx())*op.progress)
	r := image.Rect(dr.Min.X, dr.Min.Y, x1, dr.Min.Y+5)
	draw.FillRect(dst, r, color.RGBA{0, 0, 255, 128})
	draw.StrokeRect(dst, dr, op.clr, 1, -1)
}
