package main

import (
	"frame/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	initScreenWidth  = 800
	initScreenHeight = 800
)

func main() {
	ebiten.SetWindowSize(initScreenWidth, initScreenHeight)
	ebiten.SetWindowTitle("frame")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
	ui := ui.NewUI(initScreenWidth, initScreenHeight)
	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}
