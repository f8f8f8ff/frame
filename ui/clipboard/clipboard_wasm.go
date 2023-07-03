// go:build js

package clipboard

import (
	"frame/sprite"

	"github.com/hajimehoshi/ebiten/v2"
)

func copy(img *ebiten.Image) error {
	return nil
}

func paste() (*sprite.Sprite, error) {
	return nil, nil
}
