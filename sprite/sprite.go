package sprite

import (
	"frame/draw"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Sprite struct {
	Image *ebiten.Image
	Pos   image.Point
}

// returns true if point is within the bounds of the sprite
func (s Sprite) In(p image.Point) bool {
	p = p.Sub(s.Pos)
	return p.In(s.Image.Bounds())
}

// returns if sprite overlaps r
func (s Sprite) Overlaps(r image.Rectangle) bool {
	return s.Rect().Overlaps(r)
}

// returns if sprite is in r
func (s Sprite) InRect(r image.Rectangle) bool {
	return s.Rect().In(r)
}

func (s *Sprite) MoveBy(v image.Point) {
	s.Pos = s.Pos.Add(v)
}

func (s Sprite) Rect() image.Rectangle {
	return s.Image.Bounds().Add(s.Pos)
}

// draws sprite, dv is a position offset
func (s Sprite) Draw(dst *ebiten.Image, dv image.Point, alpha float64) {
	if s.Image == nil {
		return
	}
	draw.DrawImage(dst, s.Image, s.Pos.Add(dv), alpha)
}

// draw sprite with ebiten.DrawImageOptions
func (s Sprite) DrawWithOps(dst *ebiten.Image, opts *ebiten.DrawImageOptions, alpha float64) {
	copts := &colorm.DrawImageOptions{}
	copts.GeoM = opts.GeoM
	copts.GeoM.Translate(float64(s.Pos.X), float64(s.Pos.Y))
	col := colorm.ColorM{}
	col.Scale(1, 1, 1, alpha)
	colorm.DrawImage(dst, s.Image, col, copts)
}

func (s Sprite) DrawInverted(dst *ebiten.Image, dv image.Point, alpha float64) {
	draw.DrawImageInverted(dst, s.Image, s.Pos.Add(dv), alpha)
}

// resize, keep position
func (s *Sprite) Resize(newSize image.Point) {
	r := image.Rect(0, 0, newSize.X, newSize.Y)
	i := draw.ResizeImage(s.Image, r)
	s.Image = i
}

// resize, set position
func (s *Sprite) Reshape(r image.Rectangle) {
	r = r.Canon()
	newSize := r.Sub(r.Min)
	i := draw.ResizeImage(s.Image, newSize)
	s.Image = i
	s.Pos = r.Min
}

// returns a pointer to a new copy of the sprite
func (s *Sprite) Copy() *Sprite {
	return &Sprite{
		Image: ebiten.NewImageFromImage(s.Image),
		Pos:   s.Pos,
	}
}

func (s *Sprite) Crop(r image.Rectangle) *Sprite {
	im, nr := draw.CropImage(s.Image, r, s.Pos.Mul(-1))
	if im == nil {
		return nil
	}
	return &Sprite{
		Image: im,
		Pos:   nr.Min,
	}
}

// gives the sprite at position in SpriteList
func SpriteAt(sp []*Sprite, p image.Point) *Sprite {
	for i := len(sp) - 1; i >= 0; i-- {
		s := sp[i]
		if s.In(p) {
			return s
		}
	}
	return nil
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
