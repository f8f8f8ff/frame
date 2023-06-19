package ui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"

	_ "golang.org/x/image/webp"

	"github.com/hajimehoshi/ebiten/v2"
)

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
