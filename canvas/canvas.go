package canvas

import (
	"frame/draw"
	"frame/sprite"
	"image"
	"image/color"

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

func (c *Canvas) Resize(width, height int) {
	i := ebiten.NewImage(width, height)
	i.DrawImage(c.image, nil)
	c.image = i
}

func (c *Canvas) DrawSprites() {
	c.image.Fill(color.Transparent)
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

func (c Canvas) Sprites() []*sprite.Sprite {
	return c.sprites
}

func (c *Canvas) NewSpriteFromRegion(r image.Rectangle) *sprite.Sprite {
	im, r := draw.CropImage(c.image, r, image.Point{0, 0})
	if im == nil {
		return nil
	}
	return &sprite.Sprite{
		Image: im,
		Pos:   r.Min,
	}
}
