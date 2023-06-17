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
	*canvas.Canvas
	image *ebiten.Image
	err   error
	m     sync.Mutex

	operation interface{}
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

	if ui.operation == nil {
		if MouseJustPressed(ebiten.MouseButtonRight) {
			log.Println("NewMenu(MousePos())")
		}
	}

	ui.Canvas.DrawSprites()
	return nil
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// updates every frame
func (ui *UI) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	screen.DrawImage(ui.Canvas.Image(), nil)

	msg := fmt.Sprintf("%0.f", ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}
