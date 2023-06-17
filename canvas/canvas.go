package canvas

import (
	"image"
)

type Canvas struct {
	Width   int
	Height  int
	image   image.Image
	gc      *Context
	sprites []*Sprite

	cursor  image.Point
	pressed bool
}

func NewCanvas(width, height int) *Canvas {
	c := NewContext(width, height)
	c.SetRGB(1, 0, 0)
	c.Clear()
	i := c.Image()
	return &Canvas{
		Width:   width,
		Height:  height,
		image:   i,
		gc:      c,
		sprites: []*Sprite{},
		cursor:  image.Point{},
		pressed: false,
	}
}

func (c *Canvas) DrawSprites() {
	for _, s := range c.sprites {
		s.Draw(c.gc, 0, 0, 255)
	}
}

func (c Canvas) Image() image.Image {
	return c.gc.Image()
}

func (c *Canvas) AddSprite(s *Sprite) {
	c.sprites = append(c.sprites, s)
}

func (c *Canvas) RemoveSprite(s *Sprite) {
	if s == nil {
		return
	}
	index := -1
	for i, ss := range c.sprites {
		if ss == s {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.sprites = append(c.sprites[:index], c.sprites[index+1:]...)
}

func (c *Canvas) AddImage(img image.Image) {
	s := &Sprite{
		image: img,
		x:     0,
		y:     0,
	}
	c.AddSprite(s)
}
