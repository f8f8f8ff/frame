package canvas

import (
	"image"
)

type Sprite struct {
	image image.Image
	x     int
	y     int
}

func (s Sprite) In(x, y int) bool {
	p := image.Point{x - s.x, y - s.y}
	return p.In(s.image.Bounds())
}

func (s *Sprite) MoveBy(x, y int) {
	// w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()
	s.x += x
	s.y += y
}

func (s Sprite) Rect() image.Rectangle {
	return image.Rect(s.x, s.y, s.x+s.image.Bounds().Dx(), s.y+s.image.Bounds().Dy())
}

func (s Sprite) Size() image.Point {
	return s.image.Bounds().Size()
}

func (s Sprite) Draw(c *Context, dx, dy int, alpha int) {
	if alpha < 1 {
		return
	}
	if alpha < 255 {
		SetOpacity(c, alpha)
	}
	c.DrawImage(s.image, s.x+dx, s.y+dy)
}

func (s *Sprite) Resize(size image.Point) {
	r := image.Rect(0, 0, size.X, size.Y)
	i := ResizeImage(s.image, r)
	s.image = i
}

// func (s *Sprite) Draw(screen *ebiten.Image, dx, dy int, alpha float64, invert bool) {
// 	op := &colorm.DrawImageOptions{}
// 	op.GeoM.Translate(float64(s.x+dx), float64(s.y+dy))
// 	c := colorm.ColorM{}
// 	c.Scale(1, 1, 1, alpha)
// 	if invert {
// 		c.Scale(-1, -1, -1, 1)
// 		c.Translate(1, 1, 1, 0)
// 	}
// 	colorm.DrawImage(screen, s.image, c, op)
// }

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

// func (s *Sprite) Reshape(r image.Rectangle) *Sprite {
// 	oldW, oldH := s.image.Bounds().Dx(), s.image.Bounds().Dy()
// 	w, h := r.Dx(), r.Dy()
// 	if w < 1 || h < 1 {
// 		return nil
// 	}
// 	scaleX := float64(w) / float64(oldW)
// 	scaleY := float64(h) / float64(oldH)

// 	newImg := ebiten.NewImage(w, h)
// 	opt := &ebiten.DrawImageOptions{}
// 	opt.GeoM.Scale(scaleX, scaleY)
// 	newImg.DrawImage(s.image, opt)

// 	return &Sprite{
// 		image: newImg,
// 		x:     r.Min.X,
// 		y:     r.Min.Y,
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
