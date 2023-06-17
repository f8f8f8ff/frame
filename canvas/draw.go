package canvas

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"

	"github.com/fogleman/gg"
)

type Context = gg.Context

func NewContextForImage(im image.Image) *Context {
	return gg.NewContextForImage(im)
}

func NewContext(width, height int) *Context {
	return gg.NewContext(width, height)
}

func SetOpacity(c *Context, opacity int) error {
	if opacity < 0 {
		opacity = 0
	} else if opacity > 255 {
		opacity = 255
	}
	u := image.NewUniform(color.Alpha{uint8(opacity)})
	m := image.NewAlpha(image.Rect(0, 0, c.Width(), c.Height()))
	draw.Draw(m, m.Bounds(), u, image.Point{0, 0}, draw.Src)
	return c.SetMask(m)
}

func ResizeImage(src image.Image, size image.Rectangle) image.Image {
	dst := image.NewRGBA(size)
	draw.NearestNeighbor.Scale(dst, size, src, src.Bounds(), draw.Src, nil)
	return dst
}
