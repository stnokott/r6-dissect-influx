// go:build windows

package config

import (
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

func Write() {
	prefs.SetString(CONFIG_KEY_MATCH_REPLAY_FOLDER, Current.MatchReplyFolder)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_HOST, Current.InfluxDBHost)
	prefs.SetInt(CONFIG_KEY_INFLUX_DB_PORT, Current.InfluxDBPort)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_ORG, Current.InfluxDBOrg)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_BUCKET, Current.InfluxDBBucket)
	prefs.SetString(CONFIG_KEY_INFLUX_DB_TOKEN, Current.InfluxDBToken)
}
