//go:build !js && !wasm
// +build !js,!wasm

package clipboard

import (
	"bytes"
	"frame/sprite"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.design/x/clipboard"
	_ "golang.org/x/image/webp"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Println("no clipboard", err)
	}
	Enabled = true
}

func copy(img *ebiten.Image) error {
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

func paste() (*sprite.Sprite, error) {
	if !Enabled {
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
