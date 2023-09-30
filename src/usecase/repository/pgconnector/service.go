package pgconnector

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.String()), &gorm.Config{})
}

func (c Config) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		c.Host, c.User, c.Password, c.DB, c.Port)
}
