//go:build windows

// Sets up and runs the Wails application.
package main

import (
	"context"
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

// Wails hooks, will be executed in order
// targeted usecase is adding functions to the slices in build-tag-dependent source files
var onStartupFuncs, onDomReadyFuncs, onShutdownFuncs []func(context.Context)

func main() {
	// necessary for r6-dissect
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.Init()
	if err != nil {
		utils.ErrDialog(err)
		os.Exit(-1)
	}

	app := NewApp(cfg)
	onStartupFuncs = append(onStartupFuncs, app.startup)

	wailsOptions := &options.App{
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			runWailsHooks(ctx, onStartupFuncs)
		},
		OnDomReady: func(ctx context.Context) {
			runWailsHooks(ctx, onDomReadyFuncs)
		},
		OnShutdown: func(ctx context.Context) {
			runWailsHooks(ctx, onShutdownFuncs)
		},
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
			OnSuspend: func() {
				if app.roundsWatcherStop != nil {
					app.roundsWatcherStop()
				}
			},
			OnResume: func() {
				if app.roundsWatcherStop == nil {
					if errRoundWatcher := app.StartRoundWatcher(); errRoundWatcher != nil {
						utils.ErrDialog(errRoundWatcher)
					}
				}
			},
		},
		Frameless: true,
	}

	// Create application with options
	err = wails.Run(wailsOptions)

	if err != nil {
		utils.ErrDialog(err)
		os.Exit(-1)
	}
}

func runWailsHooks(ctx context.Context, hooks []func(context.Context)) {
	for _, h := range hooks {
		h(ctx)
	}
}
