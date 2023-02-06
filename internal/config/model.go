// go:build windows

package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/stnokott/r6-dissect-influx/internal/db"
)

const (
	CONFIG_KEY_MATCH_GAME_DIR   string = "game.install_dir"
	CONFIG_KEY_INFLUX_DB_HOST   string = "influx.host"
	CONFIG_KEY_INFLUX_DB_PORT   string = "influx.port"
	CONFIG_KEY_INFLUX_DB_ORG    string = "influx.org"
	CONFIG_KEY_INFLUX_DB_BUCKET string = "influx.bucket"
	CONFIG_KEY_INFLUX_DB_TOKEN  string = "influx.token"

	CONFIG_DEFAULT_INFLUX_DB_PORT int = 8086
)

type ConfigDetails struct {
	GameFolder     string
	InfluxDBHost   string
	InfluxDBPort   int
	InfluxDBOrg    string
	InfluxDBBucket string
	InfluxDBToken  string
}

func (c *ConfigDetails) InfluxURL() string {
	return "http://" + c.InfluxDBHost + ":" + strconv.Itoa(c.InfluxDBPort)
}

func (c *ConfigDetails) NewInfluxClient() *db.InfluxClient {
	return db.NewInfluxClient(db.ConnectOpts{
		URL:    c.InfluxURL(),
		Token:  c.InfluxDBToken,
		Org:    c.InfluxDBOrg,
		Bucket: c.InfluxDBBucket,
	})
}

var (
	prefs              fyne.Preferences
	Current            = new(ConfigDetails)
	bindMatchReplayDir = binding.BindString(&Current.GameFolder)
	bindInfluxHost     = binding.BindString(&Current.InfluxDBHost)
	bindInfluxPort     = binding.BindInt(&Current.InfluxDBPort)
	bindInfluxPortStr  = binding.IntToString(bindInfluxPort)
	bindInfluxOrg      = binding.BindString(&Current.InfluxDBOrg)
	bindInfluxBucket   = binding.BindString(&Current.InfluxDBBucket)
	bindInfluxToken    = binding.BindString(&Current.InfluxDBToken)
)

func Init(app fyne.App) {
	prefs = app.Preferences()
	Current.GameFolder = prefs.String(CONFIG_KEY_MATCH_GAME_DIR)
	Current.InfluxDBHost = prefs.String(CONFIG_KEY_INFLUX_DB_HOST)
	Current.InfluxDBPort = prefs.IntWithFallback(CONFIG_KEY_INFLUX_DB_PORT, CONFIG_DEFAULT_INFLUX_DB_PORT)
	Current.InfluxDBOrg = prefs.String(CONFIG_KEY_INFLUX_DB_ORG)
	Current.InfluxDBBucket = prefs.String(CONFIG_KEY_INFLUX_DB_BUCKET)
	Current.InfluxDBToken = prefs.String(CONFIG_KEY_INFLUX_DB_TOKEN)
}

func IsComplete() bool {
	return Current.GameFolder != "" &&
		Current.InfluxDBHost != "" &&
		Current.InfluxDBOrg != "" &&
		Current.InfluxDBBucket != "" &&
		Current.InfluxDBToken != ""
}

func write() {
	prefs.SetString(CONFIG_KEY_MATCH_GAME_DIR, Current.GameFolder)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_HOST, Current.InfluxDBHost)
	prefs.SetInt(CONFIG_KEY_INFLUX_DB_PORT, Current.InfluxDBPort)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_ORG, Current.InfluxDBOrg)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_BUCKET, Current.InfluxDBBucket)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_TOKEN, Current.InfluxDBToken)
}

const gameExeName string = "RainbowSix.exe"

func validateGameDir(gameDir string) (err error) {
	stats, statErr := os.Stat(gameDir)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			err = errors.New("Directory does not exist")
		} else {
			err = statErr
		}
		return
	} else if !stats.IsDir() {
		err = errors.New("Not a directory")
		return
	}

	pathToExe := path.Join(gameDir, gameExeName)
	if _, statErr = os.Stat(pathToExe); statErr != nil {
		if os.IsNotExist(statErr) {
			err = fmt.Errorf("No %s found in directory", gameExeName)
		} else {
			err = fmt.Errorf("Could not find %s in directory: %w", gameExeName, statErr)
		}
		return
	}
	return
}

var (
	regexIPv4     = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	regexHostname = regexp.MustCompile(`^(?:[0-9a-zA-Z]+\.)+[0-9a-zA-Z]{2,4}$`)
)

func validateHostAddress(s string) (err error) {
	if regexIPv4.MatchString(s) {
		return
	} else if regexHostname.MatchString(s) {
		return
	} else {
		err = errors.New("Not a valid IPv4 address or URL") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func validateInteger(s string) (err error) {
	var port int
	port, err = strconv.Atoi(s)
	if err != nil {
		err = errors.New("Not a valid integer") //lint:ignore ST1005 will be displayed in UI
	} else if port <= 0 {
		err = errors.New("Must be greater than zero") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func validateRequired(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	}
	return
}
