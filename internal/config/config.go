package config

type Config struct {
}

func (c *Config) AutodetectGameDir() (string, error) {
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

func (c *Config) ValidateGameDir(v string) error {
	return gameDirValidator.Validate(v)
}

func (c *Config) ValidateInfluxHost(v string) error {
	return influxHostValidator.Validate(v)
}

func (c *Config) ValidateInfluxPort(v string) error {
	return influxPortValidator.Validate(v)
}

func (c *Config) ValidateInfluxOrg(v string) error {
	return influxOrgValidator.Validate(v)
}

func (c *Config) ValidateInfluxBucket(v string) error {
	return influxBucketValidator.Validate(v)
}

func (c *Config) ValidateInfluxToken(v string) error {
	return influxTokenValidator.Validate(v)
}
