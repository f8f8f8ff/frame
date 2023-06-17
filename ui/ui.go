package ui

import (
	"fmt"
	"frame/canvas"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type UI struct {
	*canvas.Canvas
	image *ebiten.Image
	err   error
	m     sync.Mutex
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

	msg := fmt.Sprintf("%0.2f", ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (ui *UI) handleDroppedFiles() error {
	if err := func() error {
		ui.m.Lock()
		defer ui.m.Unlock()
		return ui.err
	}(); err != nil {
		return err
	}

	if files := ebiten.DroppedFiles(); files != nil {
		// log.Println(files)
		go func() {
			if err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fi, err := d.Info()
				if err != nil {
					return err
				}
				log.Printf("Name: %s, Size: %d, IsDir: %t, ModTime: %v", fi.Name(), fi.Size(), fi.IsDir(), fi.ModTime())

				f, err := files.Open(path)
				if err != nil {
					return err
				}

				defer func() {
					_ = f.Close()
				}()

				if fi.IsDir() {
					return nil
				}

				img, _, err := image.Decode(f)
				if err != nil {
					return err
				}

				ui.m.Lock()
				ui.Canvas.AddImage(img)
				ui.m.Unlock()

				return nil
			}); err != nil {
				ui.m.Lock()
				if ui.err == nil {
					ui.err = err
				}
				ui.m.Unlock()
			}
		}()
	}
	return nil
}
