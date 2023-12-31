package ui

import (
	"bytes"
	"frame/sprite"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"log"

	_ "golang.org/x/image/webp"

	"github.com/hajimehoshi/ebiten/v2"

	"golang.design/x/clipboard"
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

var clipboardEnabled bool

func copyClipboard(img *ebiten.Image) error {
	i := img.SubImage(img.Bounds())
	if i == nil {
		return nil
	}
	var buffer bytes.Buffer
	err := png.Encode(&buffer, i)
	if err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtImage, buffer.Bytes())
	return nil
}

func handlePaste() (*sprite.Sprite, error) {
	if !clipboardEnabled {
		return nil, nil
	}
	b := clipboard.Read(clipboard.FmtImage)
	if b == nil {
		return nil, nil
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, nil
	}
	s := &sprite.Sprite{
		Image: ebiten.NewImageFromImage(img),
	}
	return s, nil
}

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Println("no clipboard", err)
	}
	clipboardEnabled = true
}
