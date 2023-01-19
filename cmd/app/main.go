package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/root"
)

func main() {
	// TODO: display in UI
	log.Printf("%s - v%s - %s - compiled %s", constants.ProjectName, constants.Version, constants.Commit, constants.CompileTime)

	a := app.NewWithID(constants.APP_ID)

	config.Init(a)

	root.NewView(a).ShowAndRun()
}
