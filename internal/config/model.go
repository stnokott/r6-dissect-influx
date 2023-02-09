// go:build windows

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/stnokott/r6-dissect-influx/internal/db"
)

const CONFIG_DEFAULT_INFLUX_DB_PORT int = 8086

type Config struct {
	Game     GameConfigJson   `json:"game"`
	InfluxDB InfluxConfigJson `json:"influx_db"`
}

type GameConfigJson struct {
	InstallDir string `json:"install_dir"`
}

type InfluxConfigJson struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Org    string `json:"org"`
	Bucket string `json:"bucket"`
	Token  string `json:"token"`
}

func (c *Config) Read() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// config file does not exist
			setDefaults(c)
			return c.Write()
		} else {
			return err
		}
	}

	return json.Unmarshal(data, c)
}

func (c *Config) Write() error {
	if data, err := json.Marshal(c); err != nil {
		return err
	} else {
		return os.WriteFile(configPath, data, 0644)
	}
}

func (c *Config) IsComplete() bool {
	return c.Game.InstallDir != "" &&
		c.InfluxDB.Host != "" &&
		c.InfluxDB.Org != "" &&
		c.InfluxDB.Bucket != "" &&
		c.InfluxDB.Token != ""
}

func (c *Config) InfluxURL() string {
	return "http://" + c.InfluxDB.Host + ":" + strconv.Itoa(c.InfluxDB.Port)
}

func (c *Config) NewInfluxClient() *db.InfluxClient {
	return db.NewInfluxClient(db.ConnectOpts{
		URL:    c.InfluxURL(),
		Token:  c.InfluxDB.Token,
		Org:    c.InfluxDB.Org,
		Bucket: c.InfluxDB.Bucket,
	})
}

func setDefaults(target *Config) {
	target.InfluxDB.Port = CONFIG_DEFAULT_INFLUX_DB_PORT
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
