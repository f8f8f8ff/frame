package draw

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

var Font font.Face

func init() {
	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	Font = text.FaceWithLineHeight(Font, 18)
}

func BoundString(face font.Face, txt string) image.Rectangle {
	return text.BoundString(face, txt)
}

func NewTextImage(txt string, face font.Face, r image.Rectangle, padding int, fg color.Color, bg color.Color) *ebiten.Image {
	img := ebiten.NewImage(r.Dx(), r.Dy())
	img.Fill(bg)
	text.Draw(img, txt, face, padding, r.Dy()-padding, fg)
	return img
}
