//go:build windows

package main

import (
	"embed"
	"os"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/utils"

	"github.com/rs/zerolog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// TODO: move to cmd package

func main() {
	// necessary for r6-dissect
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.Init()
	if err != nil {
		utils.ErrDialog(err)
		os.Exit(-1)
	}

	app := NewApp(cfg)

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
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Frameless: true,
	})

	if err != nil {
		utils.ErrDialog(err)
		os.Exit(-1)
	}
}
