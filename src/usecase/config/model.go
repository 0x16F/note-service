package config

import (
	"notes-manager/src/usecase/repository/pgconnector"
	"notes-manager/src/usecase/repository/rsconnector"
)

type Config struct {
	Web      Web                `mapstructure:"WEB"`
	Database pgconnector.Config `mapstructure:"DATABASE"`
	Redis    rsconnector.Config `mapstructure:"REDIS"`
}

type Web struct {
	Port uint16 `mapstructure:"PORT"`
}
