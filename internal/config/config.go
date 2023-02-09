package config

import (
	"fmt"
	"os"
	"path"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

const CONFIG_FILENAME string = "config.json"

var configPath string

func Init() (*Config, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config dir: %w", err)
	}
	configDir := path.Join(userConfigDir, constants.ProjectName)
	configPath = path.Join(configDir, CONFIG_FILENAME)

	if err = os.MkdirAll(configDir, 0644); err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("could not create config dir: %w", err)
	}

	config := &Config{}
	if err = config.Read(); err != nil {
		return nil, err
	}
	return config, nil
}

type fnValidator func(string) error

type validator struct {
	fns []fnValidator
}

func newValidator(fns ...fnValidator) *validator {
	return &validator{fns: fns}
}

func (v *validator) Validate(x string) (err error) {
	for _, fn := range v.fns {
		if err = fn(x); err != nil {
			return
		}
	}
	return
}

var (
	GameDirValidator      = newValidator(validateRequired, validateGameDir)
	InfluxHostValidator   = newValidator(validateRequired, validateHostAddress)
	InfluxPortValidator   = newValidator(validateRequired, validateInteger)
	InfluxOrgValidator    = newValidator(validateRequired)
	InfluxBucketValidator = newValidator(validateRequired)
	InfluxTokenValidator  = newValidator(validateRequired)
)
