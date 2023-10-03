package migrate

import (
	"fmt"
	"notes-manager/src/usecase/repository/pgconnector"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(cfg *pgconnector.Config, down bool, migrationsPath string) error {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	m, err := migrate.New(migrationsPath, url)
	if err != nil {
		return err
	}

	if down {
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				return nil
			}

			return err
		}

		return nil
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}
