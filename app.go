package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/db"
	"github.com/stnokott/r6-dissect-influx/internal/update"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

// App struct
type App struct {
	ctx               context.Context
	config            *config.Config
	influxClient      *db.InfluxClient
	influxClientMutex sync.Mutex

	roundsWatcherStop context.CancelFunc
}

// NewApp creates a new instance of App with the provided config.
func NewApp(cfg *config.Config) *App {
	return &App{config: cfg}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// AppInfo contains information about the app itself.
type AppInfo struct {
	ProjectName string
	Version     string
	Commit      string
	GithubURL   string
}

// GetWindowTitle returns the window title from constants.
func (a *App) GetWindowTitle() string {
	return constants.WINDOW_TITLE
}

// GetAppInfo provides information about the app itself.
func (a *App) GetAppInfo() *AppInfo {
	return &AppInfo{
		ProjectName: constants.ProjectName,
		Version:     constants.Version,
		Commit:      constants.Commit,
		GithubURL:   constants.GithubURL.String(),
	}
}

// ReleaseInfo contains information about a release (of this binary).
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

// RequestLatestReleaseInfo queries information about the latest release of the app in the background
// and emits an event with the data or an error when finished.
func (a *App) RequestLatestReleaseInfo() {
	latest, err := a.getLatestRelease()
	if err == nil {
		runtime.EventsEmit(a.ctx, eventNames.LatestReleaseInfo, latest)
	} else {
		runtime.EventsEmit(a.ctx, eventNames.LatestReleaseInfoErr, err.Error())
	}
}

// StartReleaseWatcher starts a background task that regularly checks for updates.
func (a *App) StartReleaseWatcher() {
	ticker := time.NewTicker(constants.UpdateCheckInterval)
	go func() {
		for {
			a.RequestLatestReleaseInfo()
			<-ticker.C
		}
	}()
}

// StartUpdate will download the latest release
func (a *App) StartUpdate() error {
	release, err := update.GetLatestRelease()
	if err != nil {
		return err
	}
	if !release.IsNewer() {
		return errors.New("no update available")
	}
	chProgress := release.DownloadAndApply()
	go func() {
		for {
			progressInfo, ok := <-chProgress
			if !ok {
				runtime.EventsEmit(a.ctx, eventNames.UpdateProgress, "Restarting...")
				if err = utils.RestartSelf(); err != nil {
					runtime.EventsEmit(a.ctx, eventNames.UpdateErr, err)
				}
				// give time for process to launch
				time.Sleep(3 * time.Second)
				// no further event needed as app is expected to be shut down now
				return
			} else if progressInfo.Err != nil {
				runtime.EventsEmit(a.ctx, eventNames.UpdateErr, progressInfo.Err.Error())
				return
			} else {
				runtime.EventsEmit(a.ctx, eventNames.UpdateProgress, progressInfo.Task)
			}
		}
	}()
	return nil
}

// OpenGameDirDialog will open a directory selector dialog to choose the R6 game directory.
func (a *App) OpenGameDirDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Choose game directory",
	}
	return runtime.OpenDirectoryDialog(a.ctx, options)
}

// AutodetectGameDir will attempt to automatically derive the game directory from a few options, i.e. the registry.
// If the directory cannot be derived, an error will be returned.
func (*App) AutodetectGameDir() (string, error) {
	return config.GameFolderFromRegistry()
}

// ValidateGameDir checks whether the provided game directory string points to a valid
// installation directory.
func (*App) ValidateGameDir(v string) error {
	return config.GameDirValidator.Validate(v)
}

// ValidateInfluxHost performs string-only validation on the InfluxDB host (no network IO).
func (*App) ValidateInfluxHost(v string) error {
	return config.InfluxHostValidator.Validate(v)
}

// ValidateInfluxPort performs string-only validation on the InfluxDB port (no network IO).
func (*App) ValidateInfluxPort(v string) error {
	return config.InfluxPortValidator.Validate(v)
}

// ValidateInfluxOrg performs string-only validation on the InfluxDB org (no network IO).
func (*App) ValidateInfluxOrg(v string) error {
	return config.InfluxOrgValidator.Validate(v)
}

// ValidateInfluxBucket performs string-only validation on the InfluxDB bucket (no network IO).
func (*App) ValidateInfluxBucket(v string) error {
	return config.InfluxBucketValidator.Validate(v)
}

// ValidateInfluxToken performs string-only validation on the InfluxDB token (no network IO).
func (*App) ValidateInfluxToken(v string) error {
	return config.InfluxTokenValidator.Validate(v)
}

// GetConfig returns the currently persisted config.
func (a *App) GetConfig() *config.Config {
	return a.config
}

// IsConfigComplete checks if all required config fields are present in the current config.
func (a *App) IsConfigComplete() bool {
	return a.config.IsComplete()
}

// InfluxClientFromConfig creates a new InfluxDB client based on the provided config and tests if a writable connection can be established.
// If successful, the client will be saved for future use and the details of the connection are returned.
func (a *App) InfluxClientFromConfig(cfg *config.Config) (details *db.ConnectionDetails, err error) {
	if cfg == nil || !cfg.IsComplete() {
		return nil, errors.New("config incomplete, please setup first")
	}

	newClient := cfg.NewInfluxClient()
	details, err = newClient.ValidateConn(10 * time.Second)
	if err != nil {
		return
	}

	if err = cfg.Write(); err != nil {
		return
	}

	a.config = cfg
	// close old client before assigning newClient
	a.influxClientMutex.Lock()
	defer a.influxClientMutex.Unlock()
	if a.influxClient != nil {
		a.influxClient.Close()
	}
	a.influxClient = newClient
	go a.handleInfluxEvents(a.influxClient.LoopAsync())
	return
}

// InfluxClientFromSettings uses the persisted settings to create a new InfluxDB client (see InfluxClientFromConfig).
func (a *App) InfluxClientFromSettings() (details *db.ConnectionDetails, err error) {
	return a.InfluxClientFromConfig(a.config)
}

func (a *App) handleInfluxEvents(eventChan <-chan db.InfluxEvent) {
	for {
		event, ok := <-eventChan
		if !ok {
			return
		}
		runtime.EventsEmit(a.ctx, eventNames.RoundPush, event)
	}
}
