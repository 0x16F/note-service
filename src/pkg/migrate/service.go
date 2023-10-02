package migrate

import (
	"fmt"
	"notes-manager/src/usecase/repository/pgconnector"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func ApplyMigrations(cfg *pgconnector.Config, down bool, migrationsPath string) error {
	var gErr error

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	attempts := 0

	for attempts <= 4 {
		m, err := migrate.New(migrationsPath, url)
		if err != nil && attempts != 4 {
			gErr = err
			attempts += 1
			time.Sleep(3 * time.Second)
			continue
		} else {
			gErr = nil
		}

		if attempts == 4 {
			return gErr
		}

		if !down {
			if err := m.Up(); err != nil {
				logrus.Error(err)
				gErr = err
			} else {
				gErr = nil
			}
		} else {
			if err := m.Down(); err != nil {
				logrus.Error(err)
				gErr = err
			} else {
				gErr = nil
			}
		}

		break
	}

	return gErr
}
