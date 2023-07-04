package clipboard

import (
	"frame/sprite"

	"github.com/hajimehoshi/ebiten/v2"
)

var Enabled bool

func Copy(img *ebiten.Image) error {
	if !Enabled {
		return nil
	}
	return copy(img)
}

func Paste() (*sprite.Sprite, error) {
	if !Enabled {
		return nil, nil
	}
	return paste()
}
