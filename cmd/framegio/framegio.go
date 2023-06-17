package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/unit"

	"frame/uiguo"
)

var (
	initW = 800
	initH = 800
)

func main() {
	ui := uiguo.NewUI(initW, initH)

	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("frame"),
			app.Size(unit.Dp(initW), unit.Dp(initH)),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
