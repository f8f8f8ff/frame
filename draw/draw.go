package draw

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func ResizeImage(src *ebiten.Image, size image.Rectangle) *ebiten.Image {
	opts := ResizeOpts(src.Bounds(), size)
	im := ebiten.NewImage(size.Dx(), size.Dy())
	im.DrawImage(src, &opts)
	return im
}

func ReshapeOpts(src image.Rectangle, dst image.Rectangle) ebiten.DrawImageOptions {
	src = src.Canon()
	dst = dst.Canon()
	srcx, srcy := src.Dx(), src.Dy()
	dstx, dsty := dst.Dx(), dst.Dy()
	scalex := float64(srcx) / float64(dstx)
	scaley := float64(srcy) / float64(dsty)
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Scale(scalex, scaley)

	v := dst.Min.Sub(src.Min)
	trx, try := float64(v.X), float64(v.Y)
	opt.GeoM.Translate(trx, try)
	return opt
}

func ResizeOpts(src image.Rectangle, dst image.Rectangle) ebiten.DrawImageOptions {
	src = src.Canon()
	dst = dst.Canon()
	src = src.Sub(src.Min)
	dst = dst.Sub(dst.Min)
	return ReshapeOpts(src, dst)
}

func CropImage(src *ebiten.Image, r image.Rectangle, offset image.Point) (*ebiten.Image, image.Rectangle) {
	r = r.Canon()
	r = r.Add(offset)
	if r.Dx() < 1 || r.Dy() < 1 {
		return nil, r
	}
	if !r.Overlaps(src.Bounds()) {
		return nil, r
	}
	intr := r.Intersect(src.Bounds())
	im := src.SubImage(intr)
	intr = intr.Sub(offset)
	return ebiten.NewImageFromImage(im), intr
}

func DrawImage(dst, src *ebiten.Image, pos image.Point, alpha float64) {
	opts := &colorm.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	col := colorm.ColorM{}
	col.Scale(1, 1, 1, alpha)
	colorm.DrawImage(dst, src, col, opts)
}

func DrawImageInverted(dst, src *ebiten.Image, pos image.Point, alpha float64) {
	opts := &colorm.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	col := colorm.ColorM{}
	col.Scale(1, 1, 1, alpha)
	col.Scale(-1, -1, -1, 1)
	col.Translate(1, 1, 1, 0)
	colorm.DrawImage(dst, src, col, opts)
}

// outline
func StrokeRect(dst *ebiten.Image, rect image.Rectangle, clr color.Color, strokeWidth, offset float32) {
	x := float32(rect.Min.X) - offset/2
	y := float32(rect.Min.Y) - offset/2
	w := float32(rect.Dx()) + offset
	h := float32(rect.Dy()) + offset
	vector.StrokeRect(dst, x, y, w, h, strokeWidth, clr, false)
}
