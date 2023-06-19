package ui

import (
	"fmt"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"frame/canvas"
)

type UI struct {
	Width  int
	Height int

	*canvas.Canvas
	image *ebiten.Image
	err   error
	m     sync.Mutex

	operations []interface{}
	LockOrder  bool
}

func NewUI(w, h int) *UI {
	c := canvas.NewCanvas(w, h)
	i := ebiten.NewImage(w, h)
	return &UI{
		Canvas: c,
		image:  i,
	}
}

// updates on ticks
func (ui *UI) Update() error {
	err := ui.handleDroppedFiles()
	if err != nil {
		return err
	}

	if len(ui.operations) == 0 {
		if MouseJustPressed(ebiten.MouseButtonRight) {
			ui.addOperation(MainMenu())
		}
		if MouseJustPressed(ebiten.MouseButtonLeft) {
			ui.addOperation(&DragOp{})
		}
	}
	ui.HandleOperations()

	ui.Canvas.DrawSprites()
	return nil
}

func (ui *UI) HandleOperations() (err error) {
	for _, ope := range ui.operations {
		if ope == nil {
			ui.removeOperation(ope)
			continue
		}
		switch op := ope.(type) {
		case *Menu:
			if o, done := op.Update(); done {
				ui.removeOperation(op)
				ui.addOperation(o)
			}
		case Operation:
			if CancelInput() {
				ui.operations = []interface{}{}
			}
			if done, e := op.Update(ui); done {
				ui.removeOperation(op)
				err = e
			}
		default:
			log.Printf("unhandled operation: %T %v", op, op)
			ui.removeOperation(op)
		}
	}
	return err
}

func (ui *UI) Layout(newWidth, newHeight int) (int, int) {
	if newWidth != ui.Width || newHeight != ui.Height {
		ui.Width = newWidth
		ui.Height = newHeight
		ui.Canvas.Resize(ui.Width, ui.Height)
	}
	return ui.Width, ui.Height
}

// updates every frame
func (ui *UI) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	screen.DrawImage(ui.Canvas.Image(), nil)

	for _, ope := range ui.operations {
		switch op := ope.(type) {
		case Drawable:
			op.Draw(screen)
		}
	}

	msg := fmt.Sprintf("%0.f\n", ebiten.ActualFPS())
	if len(ui.operations) > 0 {
		for _, o := range ui.operations {
			msg += fmt.Sprintf("%T %v\n", o, o)
		}
	}
	ebitenutil.DebugPrint(screen, msg)
}

func (ui *UI) addOperation(op interface{}) {
	ui.operations = append(ui.operations, op)
}

func (ui *UI) removeOperation(op interface{}) {
	index := -1
	for i, o := range ui.operations {
		if o == op {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	ui.operations = append(ui.operations[:index], ui.operations[index+1:]...)
}
