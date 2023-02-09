package main

import (
	"context"
	"time"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/db"
)

// App struct
type App struct {
	ctx          context.Context
	config       *config.Config
	influxClient *db.InfluxClient
}

func NewApp(cfg *config.Config) *App {
	return &App{config: cfg}
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

func (a *App) GetVersion() *BuildInfo {
	return &BuildInfo{
		Version: constants.Version,
		Commit:  constants.Commit,
	}
}

func (_ *App) AutodetectGameDir() (string, error) {
	return config.GameFolderFromRegistry()
}

func (_ *App) ValidateGameDir(v string) error {
	return config.GameDirValidator.Validate(v)
}

func (_ *App) ValidateInfluxHost(v string) error {
	return config.InfluxHostValidator.Validate(v)
}

func (_ *App) ValidateInfluxPort(v string) error {
	return config.InfluxPortValidator.Validate(v)
}

func (_ *App) ValidateInfluxOrg(v string) error {
	return config.InfluxOrgValidator.Validate(v)
}

func (_ *App) ValidateInfluxBucket(v string) error {
	return config.InfluxBucketValidator.Validate(v)
}

func (_ *App) ValidateInfluxToken(v string) error {
	return config.InfluxTokenValidator.Validate(v)
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) SaveAndValidateConfig(cfg *config.Config) (details *db.ConnectionDetails, err error) {
	newClient := cfg.NewInfluxClient()
	details, err = newClient.ValidateConn(10 * time.Second)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			newClient.Close()
		}
	}()

	if err = a.config.Write(); err != nil {
		return
	}

	a.config = cfg
	a.influxClient = newClient
	return
}
