package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/root"
)

func main() {
	a := app.NewWithID(constants.APP_ID)

	config.Init(a)

	root.NewView(a).ShowAndRun()
}
