package main

import (
	"context"
	"errors"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/db"
	"github.com/stnokott/r6-dissect-influx/internal/update"
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

type AppInfo struct {
	ProjectName string
	Version     string
	Commit      string
	GithubURL   string
}

func (a *App) GetWindowTitle() string {
	return constants.WINDOW_TITLE
}

func (a *App) GetAppInfo() *AppInfo {
	return &AppInfo{
		ProjectName: constants.ProjectName,
		Version:     constants.Version,
		Commit:      constants.Commit,
		GithubURL:   constants.GithubURL.String(),
	}
}

type ReleaseInfo struct {
	IsNewer     bool
	Version     string
	Commitish   string
	PublishedAt time.Time
	Changelog   string
}

func (*App) getLatestRelease() (*ReleaseInfo, error) {
	release, err := update.GetLatestRelease()
	if err != nil {
		return nil, err
	}
	if len(release.Commitish) > 7 {
		release.Commitish = release.Commitish[:7]
	}
	return &ReleaseInfo{
		IsNewer:     release.IsNewer(),
		Version:     release.SemVer.String(),
		Commitish:   release.Commitish,
		PublishedAt: release.PublishedAt,
		Changelog:   release.Body,
	}, nil
}

func (a *App) RequestLatestReleaseInfo() {
	latest, err := a.getLatestRelease()
	if err == nil {
		runtime.EventsEmit(a.ctx, eventNames.LatestReleaseInfo, latest)
	} else {
		runtime.EventsEmit(a.ctx, eventNames.LatestReleaseInfoErr, err)
	}
}

func (a *App) StartReleaseWatcher() {
	ticker := time.NewTicker(constants.UpdateCheckInterval)
	go func() {
		for {
			a.RequestLatestReleaseInfo()
			<-ticker.C
		}
	}()
}

type UpdateProgress struct {
	Task     string
	Complete bool
}

func (a *App) StartUpdate() error {
	release, err := update.GetLatestRelease()
	if err != nil {
		return err
	}
	chProgress := release.DownloadAndApply()
	go func() {
		for {
			progressInfo, ok := <-chProgress
			if !ok {
				runtime.EventsEmit(a.ctx, eventNames.UpdateProgress, UpdateProgress{Complete: true})
				return
			} else if progressInfo.Err != nil {
				runtime.EventsEmit(a.ctx, eventNames.UpdateErr, progressInfo.Err.Error())
				return
			} else {
				runtime.EventsEmit(a.ctx, eventNames.UpdateProgress, UpdateProgress{Task: progressInfo.Task})
			}
		}
	}()
	return nil
}

func (a *App) OpenGameDirDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Choose game directory",
	}
	return runtime.OpenDirectoryDialog(a.ctx, options)
}

func (*App) AutodetectGameDir() (string, error) {
	return config.GameFolderFromRegistry()
}

func (*App) ValidateGameDir(v string) error {
	return config.GameDirValidator.Validate(v)
}

func (*App) ValidateInfluxHost(v string) error {
	return config.InfluxHostValidator.Validate(v)
}

func (*App) ValidateInfluxPort(v string) error {
	return config.InfluxPortValidator.Validate(v)
}

func (*App) ValidateInfluxOrg(v string) error {
	return config.InfluxOrgValidator.Validate(v)
}

func (*App) ValidateInfluxBucket(v string) error {
	return config.InfluxBucketValidator.Validate(v)
}

func (*App) ValidateInfluxToken(v string) error {
	return config.InfluxTokenValidator.Validate(v)
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) IsConfigComplete() bool {
	return a.config.IsComplete()
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

	if err = cfg.Write(); err != nil {
		return
	}

	a.config = cfg
	a.influxClient = newClient
	return
}

func (a *App) Connect() (details *db.ConnectionDetails, err error) {
	if a.config == nil || !a.config.IsComplete() {
		return nil, errors.New("config incomplete, please setup first")
	}
	if a.influxClient == nil {
		a.influxClient = a.config.NewInfluxClient()
	}
	return a.influxClient.ValidateConn(10 * time.Second)
}
