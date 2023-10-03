package main

import (
	"notes-manager/src/controller/web"
	"notes-manager/src/pkg/fsscanner"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/config"
	"notes-manager/src/usecase/repository"
	"notes-manager/src/usecase/repository/pgconnector"
	"notes-manager/src/usecase/repository/rsconnector"
	"os"

	"github.com/sirupsen/logrus"
)

// @title Notes Service API
// @version 1.0
// @description Simple server to demonstrate some features

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

	// getting current path to the current folder
	wd, _ := os.Getwd()

	// getting path to the migrations
	migrationsPath, err := fsscanner.FindDirectory(wd, "migrations")
	if err != nil {
		logrus.Fatal("failed to find migrations folder")
	}

	// apply migrations
	if err := migrate.ApplyMigrations(&cfg.Database, false, "file://"+migrationsPath); err != nil {
		logrus.Fatalf("failed to apply migrations, %v", err)
	}

	// init && start web
	if err := web.New(repo).Start(cfg.Web.Port); err != nil {
		logrus.Fatalf("failed to start web server, %v", err)
	}
}
