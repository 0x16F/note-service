package dcontainer

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgres(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16.0-alpine3.18",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "note-service",
			"POSTGRES_DB":       "note_service",
		},
	}

	gcr := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	return testcontainers.GenericContainer(ctx, gcr)
}
