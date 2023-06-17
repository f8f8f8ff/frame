package uiguo

import (
	"fmt"
	"frame/canvas"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"gioui.org/app"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type UI struct {
	*canvas.Canvas
}

func NewUI(w, h int) *UI {
	c := canvas.NewCanvas(w, h)
	return &UI{
		Canvas: c,
	}
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			fmt.Println("frame")
			gtx := layout.NewContext(&ops, e)
			stop := ui.HandleKeyboard(gtx, w)
			if stop {
				return nil
			}
			err := ui.HandleDragDrop(gtx, w)
			if err != nil {
				return err
			}
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.NW.Layout(gtx,
		CanvasWidget{
			Canvas: ui.Canvas,
		}.Layout,
	)
}

type CanvasWidget struct {
	*canvas.Canvas
}

func (c CanvasWidget) Layout(gtx layout.Context) layout.Dimensions {
	// update the canvas
	c.Canvas.DrawSprites()

	size := image.Point{c.Width, c.Height}
	gtx.Constraints = layout.Exact(size)

	imageOp := paint.NewImageOp(c.Image())
	imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (ui *UI) HandleDragDrop(gtx layout.Context, w *app.Window) error {
	fmt.Println("HandleDragDrop")
	transfer.TargetOp{
		Tag:  w,
		Type: "image/jpeg",
	}.Add(gtx.Ops)
	transfer.TargetOp{
		Tag:  w,
		Type: "image/png",
	}.Add(gtx.Ops)

	// add webp

	for _, e := range gtx.Events(w) {
		switch e := e.(type) {
		case transfer.DataEvent:
			file := e.Open()
			img, err := loadImage(file)
			if err != nil {
				return err
			}
			ui.Canvas.AddImage(img)
		}
	}
	return nil
}

func (ui *UI) HandleKeyboard(gtx layout.Context, w *app.Window) bool {
	clipboard.ReadOp{
		Tag: w,
	}.Add(gtx.Ops)
	key.InputOp{
		Tag:  w,
		Keys: key.NameEscape,
	}.Add(gtx.Ops)
	for _, event := range gtx.Events(w) {
		switch event := event.(type) {
		case key.Event:
			if event.Name == key.NameEscape {
				return true
			}
		case clipboard.Event:
			fmt.Println(event.Text)
		}

	}
	return false
}

func loadImage(file io.ReadCloser) (image.Image, error) {
	img, _, err := image.Decode(file)
	return img, err
}
