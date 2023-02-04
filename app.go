package main

import (
	"context"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type BuildInfo struct {
	Version string
	Commit  string
}

func (a *App) GetWindowTitle() string {
	return constants.WINDOW_TITLE
}

// Greet returns a greeting for the given name
func (a *App) GetVersion() *BuildInfo {
	return &BuildInfo{
		Version: constants.Version,
		Commit:  constants.Commit,
	}
}
