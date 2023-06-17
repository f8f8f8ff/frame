package canvas

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
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

func DrawImage(dst, src *ebiten.Image, pos image.Point, alpha float64) {
	opts := &colorm.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	col := colorm.ColorM{}
	col.Scale(1, 1, 1, alpha)
	colorm.DrawImage(dst, src, col, opts)
}
