package main

import (
	"frame/uiebiten"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	initScreenWidth  = 800
	initScreenHeight = 800
)

func main() {
	ebiten.SetWindowSize(initScreenWidth, initScreenHeight)
	ebiten.SetWindowTitle("title")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(false)
	ui := uiebiten.NewUI(initScreenWidth, initScreenHeight)
	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}
