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
	Sprites sprite.SpriteList

	cursor  image.Point
	pressed bool
}

func NewCanvas(width, height int) *Canvas {
	i := ebiten.NewImage(width, height)
	return &Canvas{
		Width:   width,
		Height:  height,
		image:   i,
		Sprites: []*sprite.Sprite{},
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
	for i := len(c.Sprites) - 1; i >= 0; i-- {
		s := c.Sprites[i]
		s.Draw(c.image, image.Point{0, 0}, 1)
	}
}

func (c Canvas) Image() *ebiten.Image {
	return c.image
}

func (c *Canvas) AddSprite(s *sprite.Sprite) {
	c.Sprites = append(sprite.SpriteList{s}, c.Sprites...)
}

func (c *Canvas) RemoveSprite(s *sprite.Sprite) {
	if s == nil {
		return
	}
	index := -1
	for i, ss := range c.Sprites {
		if ss == s {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.Sprites = append(c.Sprites[:index], c.Sprites[index+1:]...)
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
	return sprite.SpriteAt(c.Sprites, p)
}

func (c Canvas) GetSprites() []*sprite.Sprite {
	return c.Sprites
}

func (c *Canvas) Reorder(command sprite.ReorderCommand, s *sprite.Sprite) {
	newList := c.Sprites.Reorder(command, s)
	c.Sprites = newList
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
