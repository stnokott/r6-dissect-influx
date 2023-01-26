//go:build windows && cgo

//go:generate goversioninfo -64

package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/rs/zerolog"
	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/root"
)

func main() {
	// necessary for r6-dissect
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	a := app.NewWithID(constants.APP_ID)
	a.SetIcon(constants.APP_ICON)

	config.Init(a)

	root.NewView(a).ShowAndRun()
}
