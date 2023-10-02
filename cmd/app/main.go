package main

import (
	"notes-manager/src/controller/web"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/config"
	"notes-manager/src/usecase/repository"
	"notes-manager/src/usecase/repository/pgconnector"
	"notes-manager/src/usecase/repository/rsconnector"

	"github.com/sirupsen/logrus"
)

func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		logrus.Fatalf("failed to read the config, %v", err)
	}

	// connecting to the database
	db, err := pgconnector.Connect(&cfg.Database)
	if err != nil {
		logrus.Fatalf("failed to connect to the database, %v", err)
	}

	// connecting to the redis
	client, err := rsconnector.Connect(&cfg.Redis)
	if err != nil {
		logrus.Fatalf("failed to connect to the redis, %v", err)
	}

	// init repo
	repo := repository.New(db, client)

	// apply migrations
	if err := migrate.ApplyMigrations(&cfg.Database, false, "file://migrations"); err != nil {
		logrus.Fatalf("failed to apply migrations, %v", err)
	}

	// init && start web
	if err := web.New(repo).Start(cfg.Web.Port); err != nil {
		logrus.Fatalf("failed to start web server, %v", err)
	}
}
