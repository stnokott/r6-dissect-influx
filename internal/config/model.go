// go:build windows

package config

import (
	"errors"
	"os"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const (
	CONFIG_KEY_MATCH_REPLAY_FOLDER string = "game.replay_folder"
	CONFIG_KEY_INFLUX_DB_HOST      string = "influx.host"
	CONFIG_KEY_INFLUX_DB_PORT      string = "influx.port"
	CONFIG_KEY_INFLUX_DB_ORG       string = "influx.org"
	CONFIG_KEY_INFLUX_DB_BUCKET    string = "influx.bucket"
	CONFIG_KEY_INFLUX_DB_TOKEN     string = "influx.token"

	CONFIG_DEFAULT_INFLUX_DB_PORT int = 8086
)

type Config struct {
	MatchReplyFolder string
	InfluxDBHost     string
	InfluxDBPort     int
	InfluxDBOrg      string
	InfluxDBBucket   string
	InfluxDBToken    string
}

var (
	prefs              fyne.Preferences
	Current            = new(Config)
	bindMatchReplayDir = binding.BindString(&Current.MatchReplyFolder)
	bindInfluxHost     = binding.BindString(&Current.InfluxDBHost)
	bindInfluxPort     = binding.BindInt(&Current.InfluxDBPort)
	bindInfluxPortStr  = binding.IntToString(bindInfluxPort)
	bindInfluxOrg      = binding.BindString(&Current.InfluxDBOrg)
	bindInfluxBucket   = binding.BindString(&Current.InfluxDBBucket)
	bindInfluxToken    = binding.BindString(&Current.InfluxDBToken)
)

func Init(app fyne.App) {
	prefs = app.Preferences()
	Current.MatchReplyFolder = prefs.StringWithFallback(CONFIG_KEY_MATCH_REPLAY_FOLDER, "")
	Current.InfluxDBHost = prefs.StringWithFallback(CONFIG_KEY_INFLUX_DB_HOST, "")
	Current.InfluxDBPort = prefs.IntWithFallback(CONFIG_KEY_INFLUX_DB_PORT, CONFIG_DEFAULT_INFLUX_DB_PORT)
	Current.InfluxDBOrg = prefs.StringWithFallback(CONFIG_KEY_INFLUX_DB_ORG, "")
	Current.InfluxDBBucket = prefs.StringWithFallback(CONFIG_KEY_INFLUX_DB_BUCKET, "")
	Current.InfluxDBToken = prefs.StringWithFallback(CONFIG_KEY_INFLUX_DB_TOKEN, "")
}

func write() {
	prefs.SetString(CONFIG_KEY_MATCH_REPLAY_FOLDER, Current.MatchReplyFolder)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_HOST, Current.InfluxDBHost)
	prefs.SetInt(CONFIG_KEY_INFLUX_DB_PORT, Current.InfluxDBPort)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_ORG, Current.InfluxDBOrg)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_BUCKET, Current.InfluxDBBucket)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_TOKEN, Current.InfluxDBToken)
}

// TODO: warn if no matches found or none could be parsed
func directoryValidator(s string) (err error) {
	stats, statErr := os.Stat(s)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			err = errors.New("Does not exist") //lint:ignore ST1005 will be displayed in UI
		} else {
			err = statErr
		}
	} else if !stats.IsDir() {
		err = errors.New("Not a directory") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

var (
	regexIPv4     = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	regexHostname = regexp.MustCompile(`^(?:[0-9a-zA-Z]+\.)+[0-9a-zA-Z]{2,4}$`)
)

func hostAddressValidator(s string) (err error) {
	if regexIPv4.MatchString(s) {
		return
	} else if regexHostname.MatchString(s) {
		return
	} else {
		err = errors.New("Not a valid IPv4 address or URL") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func integerValidator(s string) (err error) {
	var port int
	port, err = strconv.Atoi(s)
	if err != nil {
		err = errors.New("Not a valid integer") //lint:ignore ST1005 will be displayed in UI
	} else if port <= 0 {
		err = errors.New("Must be greater than zero") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func requiredValidator(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	}
	return
}
