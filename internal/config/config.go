package config

type Config struct {
	Logger struct {
		Level  string `mapstructure:"level" env:"LOG_LEVEL" default:"info"`
		Format string `mapstructure:"format" env:"LOG_FORMAT" default:"json"`
		Env    string `mapstructure:"env" env:"APP_ENV" default:"development"`
	}
}
