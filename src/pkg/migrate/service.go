package migrate

import (
	"fmt"
	"notes-manager/src/usecase/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func ApplyMigrations(down bool) error {
	var gErr error

	// Инициализируем конфиг
	cfg, err := config.New()
	if err != nil {
		logrus.Fatal(err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DB)

	attempts := 0

	for attempts <= 4 {
		m, err := migrate.New("file://migrations", url)
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
