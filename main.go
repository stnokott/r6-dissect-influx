//go:build windows

//go:generate goversioninfo -64

package main

import (
	"embed"
	"os"

	"github.com/tawesoft/golib/v2/dialog"

	"github.com/rs/zerolog"
	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// necessary for r6-dissect
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_ = dialog.Init()

	app := NewApp()
	cfg, err := config.New()
	if err != nil {
		_ = dialog.Error(err.Error())
		os.Exit(-1)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			cfg,
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Frameless: true,
	})

	if err != nil {
		_ = dialog.Error(err.Error())
		os.Exit(-1)
	}
}
