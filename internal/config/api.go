package config

import (
	"fmt"
	"os"
	"path"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

const CONFIG_FILENAME string = "config.json"

type API struct {
	configPath string
	config     *ConfigJSON
}

func New() (*API, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config dir: %w", err)
	}
	configDir := path.Join(userConfigDir, constants.ProjectName)
	configPath := path.Join(configDir, CONFIG_FILENAME)

	if err = os.MkdirAll(configDir, 0644); err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("could not create config dir: %w", err)
	}

	config := &ConfigJSON{}
	if err = read(configPath, config); err != nil {
		return nil, err
	}
	return &API{
		configPath: configPath,
		config:     config,
	}, nil
}

func (_ *API) AutodetectGameDir() (string, error) {
	return gameFolderFromRegistry()
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
	gameDirValidator      = newValidator(validateRequired, validateGameDir)
	influxHostValidator   = newValidator(validateRequired, validateHostAddress)
	influxPortValidator   = newValidator(validateRequired, validateInteger)
	influxOrgValidator    = newValidator(validateRequired)
	influxBucketValidator = newValidator(validateRequired)
	influxTokenValidator  = newValidator(validateRequired)
)

func (_ *API) ValidateGameDir(v string) error {
	return gameDirValidator.Validate(v)
}

func (_ *API) ValidateInfluxHost(v string) error {
	return influxHostValidator.Validate(v)
}

func (_ *API) ValidateInfluxPort(v string) error {
	return influxPortValidator.Validate(v)
}

func (_ *API) ValidateInfluxOrg(v string) error {
	return influxOrgValidator.Validate(v)
}

func (_ *API) ValidateInfluxBucket(v string) error {
	return influxBucketValidator.Validate(v)
}

func (_ *API) ValidateInfluxToken(v string) error {
	return influxTokenValidator.Validate(v)
}

func (a *API) GetConfig() *ConfigJSON {
	return a.config
}

func (a *API) SaveConfig(c *ConfigJSON) error {
	a.config = c
	return write(c, a.configPath)
}
