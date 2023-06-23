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
	if r.Dx() == 0 || r.Dy() == 0 {
		return nil
	}
	img := ebiten.NewImage(r.Dx(), r.Dy())
	img.Fill(bg)
	text.Draw(img, txt, face, padding, r.Dy()-padding, fg)
	return img
}

func TextLineImage(txt string, face font.Face, height, padding int, fg color.Color, bg color.Color) *ebiten.Image {
	if height < 1 {
		return nil
	}
	width := BoundString(face, txt).Dx() + 2*padding
	r := image.Rect(0, 0, width, height)
	return NewTextImage(txt, face, r, padding, fg, bg)
}

func NewTextBlockImage(txt string, face font.Face, padding int, fg color.Color, bg color.Color) *ebiten.Image {
	r := BoundString(face, txt)
	// r = r.Sub(r.Min)
	r = image.Rect(0, 0, r.Dx()+2*padding, r.Dy()+2*padding)
	img := ebiten.NewImage(r.Dx(), r.Dy())
	img.Fill(bg)
	text.Draw(img, txt, face, padding, face.Metrics().Height.Ceil(), fg)
	return img
}
