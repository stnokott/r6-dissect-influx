// go:build windows

package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/stnokott/r6-dissect-influx/constants"
	"golang.org/x/sys/windows/registry"
)

const (
	CONFIG_KEY_MATCH_REPLAY_FOLDER string = "game.replay_folder"
	CONFIG_KEY_INFLUX_DB_HOST      string = "influx.host"
	CONFIG_KEY_INFLUX_DB_PORT      string = "influx.port"

	CONFIG_DEFAULT_INFLUX_DB_PORT int = 8086
)

type Config struct {
	MatchReplyFolder string
	InfluxDBHost     string
	InfluxDBPort     int
	InfluxDBBucket   string
}

func Load(configFilePath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFilePath)

	viper.SetDefault(CONFIG_KEY_INFLUX_DB_PORT, CONFIG_DEFAULT_INFLUX_DB_PORT)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	return &Config{
		MatchReplyFolder: viper.GetString(CONFIG_KEY_MATCH_REPLAY_FOLDER),
		InfluxDBHost:     viper.GetString(CONFIG_KEY_INFLUX_DB_HOST),
		InfluxDBPort:     viper.GetInt(CONFIG_KEY_INFLUX_DB_PORT),
	}, nil
}

const gameFolderRegistryKey string = `SOFTWARE\WOW6432Node\Ubisoft\Launcher\Installs\635`

func matchReplayFolderFromRegistry() (result string, err error) {
	var key registry.Key
	key, err = registry.OpenKey(registry.LOCAL_MACHINE, gameFolderRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer func() {
		errInner := key.Close()
		if errInner != nil && err == nil {
			err = errInner
		}
	}()
	var gameDir string
	gameDir, _, err = key.GetStringValue("InstallDir")
	if err != nil {
		return
	}

	_, err = os.Stat(gameDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(`game directory "%s" found in Registry, but does not exist`, gameDir)
		} else {
			err = fmt.Errorf(`game directory "%s" found in Registry, but could not read folder: %w`, gameDir, err)
		}
		return
	}

	result = filepath.Join(gameDir, constants.DEFAULT_MATCH_REPLAY_FOLDER_NAME)
	var folderInfo fs.FileInfo
	folderInfo, err = os.Stat(result)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(`folder "%s" in game directory "%s" not found`, constants.DEFAULT_MATCH_REPLAY_FOLDER_NAME, gameDir)
		} else {
			err = fmt.Errorf(`could not read folder "%s"`, result)
		}
		return
	} else if !folderInfo.IsDir() {
		err = fmt.Errorf(`"%s" is not a folder`, result)
	}

	return
}
