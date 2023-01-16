// go:build windows

package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	CONFIG_KEY_MATCH_REPLAY_FOLDER string = "game.replay_folder"
	CONFIG_KEY_INFLUX_DB_HOST      string = "influx.host"
	CONFIG_KEY_INFLUX_DB_PORT      string = "influx.port"
	CONFIG_KEY_INFLUX_DB_BUCKET    string = "influx.bucket"

	CONFIG_DEFAULT_INFLUX_DB_PORT int = 8086
)

type Config struct {
	MatchReplyFolder string
	InfluxDBHost     string
	InfluxDBPort     int
	InfluxDBBucket   string
}

var CurrentConfig *Config = nil

const configFilePath string = "."

func configureViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFilePath)

	viper.SetDefault(CONFIG_KEY_INFLUX_DB_PORT, CONFIG_DEFAULT_INFLUX_DB_PORT)
}

func Init() (err error) {
	configureViper()

	if err = viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); isNotFound {
			// config file not found, write default config
			if errWrite := viper.SafeWriteConfig(); errWrite != nil {
				err = fmt.Errorf("could not write default config: %w", errWrite)
				return
			} else {
				log.Println("no config file found, wrote default config")
				// we successfully wrote a default config, so no need to return an error
				err = nil
			}
		} else {
			err = fmt.Errorf("could not read config: %w", err)
			return
		}
	}
	log.Println("config read")
	CurrentConfig = &Config{
		MatchReplyFolder: viper.GetString(CONFIG_KEY_MATCH_REPLAY_FOLDER),
		InfluxDBHost:     viper.GetString(CONFIG_KEY_INFLUX_DB_HOST),
		InfluxDBPort:     viper.GetInt(CONFIG_KEY_INFLUX_DB_PORT),
		InfluxDBBucket:   viper.GetString(CONFIG_KEY_INFLUX_DB_BUCKET),
	}
	return
}

func Write() (err error) {
	viper.Set(CONFIG_KEY_MATCH_REPLAY_FOLDER, CurrentConfig.MatchReplyFolder)
	viper.Set(CONFIG_KEY_INFLUX_DB_HOST, CurrentConfig.InfluxDBHost)
	viper.Set(CONFIG_KEY_INFLUX_DB_PORT, CurrentConfig.InfluxDBPort)
	viper.Set(CONFIG_KEY_INFLUX_DB_BUCKET, CurrentConfig.InfluxDBBucket)

	err = viper.WriteConfig()
	if err == nil {
		log.Println("config written")
	}
	return
}
