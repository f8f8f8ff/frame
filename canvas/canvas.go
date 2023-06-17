package canvas

import (
	"frame/sprite"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Canvas struct {
	Width   int
	Height  int
	image   *ebiten.Image
	sprites []*sprite.Sprite

	cursor  image.Point
	pressed bool
}

func NewCanvas(width, height int) *Canvas {
	i := ebiten.NewImage(width, height)
	return &Canvas{
		Width:   width,
		Height:  height,
		image:   i,
		sprites: []*sprite.Sprite{},
		cursor:  image.Point{},
		pressed: false,
	}
}

func (c *Canvas) DrawSprites() {
	for _, s := range c.sprites {
		s.Draw(c.image, image.Point{0, 0}, 1)
	}
}

func (c Canvas) Image() *ebiten.Image {
	return c.image
}

func (c *Canvas) AddSprite(s *sprite.Sprite) {
	c.sprites = append(c.sprites, s)
}

func (c *Canvas) RemoveSprite(s *sprite.Sprite) {
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
	i := ebiten.NewImageFromImage(img)
	s := &sprite.Sprite{
		Image: i,
		Pos:   image.Point{0, 0},
	}
	c.AddSprite(s)
}

func (c *Canvas) SpriteAt(p image.Point) *sprite.Sprite {
	return sprite.SpriteAt(c.sprites, p)
}
