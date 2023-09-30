package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

func New() (*Config, error) {
	viper := viper.NewWithOptions(
		viper.KeyDelimiter("_"),
	)

	viper.AddConfigPath(os.Getenv("PWD") + "/configs")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Join(err, errors.New("failed to read config"))
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Join(err, errors.New("failed to unmarshal config"))
	}

	return &config, nil
}
