package canvas

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image *ebiten.Image
	pos   image.Point
}

func (s Sprite) In(p image.Point) bool {
	p = p.Sub(s.pos)
	return p.In(s.image.Bounds())
}

func (s *Sprite) MoveBy(v image.Point) {
	s.pos = s.pos.Add(v)
}

func (s Sprite) Rect() image.Rectangle {
	return s.image.Bounds().Add(s.pos)
}

func (s Sprite) Draw(dst *ebiten.Image, dv image.Point, alpha float64) {
	DrawImage(dst, s.image, s.pos.Add(dv), alpha)
}

func (s *Sprite) Resize(size image.Point) {
	r := image.Rect(0, 0, size.X, size.Y)
	i := ResizeImage(s.image, r)
	s.image = i
}

func (s *Sprite) Reshape(r image.Rectangle) {

}

// func (s *Sprite) Crop(r image.Rectangle) (newSprite *Sprite) {
// 	x := r.Min.X - s.x
// 	y := r.Min.Y - s.y
// 	dx := r.Dx()
// 	dy := r.Dy()
// 	if dx < 2 || dy < 2 {
// 		return nil
// 	}
// 	dx += x
// 	dy += y
// 	cropRect := image.Rect(x, y, dx, dy)
// 	if !cropRect.Overlaps(s.image.Bounds()) {
// 		return nil
// 	}
// 	newRect := r.Intersect(s.Rect())

// 	croppedImage := s.image.SubImage(cropRect)

// 	return &Sprite{
// 		image: ebiten.NewImageFromImage(croppedImage),
// 		x:     newRect.Min.X,
// 		y:     newRect.Min.Y,
// 	}
// }

// func (s *Sprite) Copy() Sprite {
// 	return Sprite{
// 		image: ebiten.NewImageFromImage(s.image),
// 		x:     s.x,
// 		y:     s.y,
// 	}
// }

// func (s *Sprite) Outline(screen *ebiten.Image, strokeWidth float32, clr color.Color, outer bool) {
// 	w := s.image.Bounds().Dx()
// 	h := s.image.Bounds().Dy()
// 	x := s.x
// 	y := s.y
// 	if outer {
// 		StrokeRectOuter(screen, x, y, w, h, strokeWidth, clr, false)
// 		return
// 	}
// 	StrokeRectInner(screen, x, y, w, h, strokeWidth, clr, false)
// }
