package migration

import (
	"context"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/infiniteloopcloud/log"
)

type PreHookMigrationParams struct {
	SkipMigration    string
	ConnectionString string
	PostgresDialect  bool
	Files            []embed.FS
}

func PreHookMigration(params PreHookMigrationParams) func(ctx context.Context) {
	return func(ctx context.Context) {
		if params.SkipMigration != "true" {
			m, err := Get(
				ctx,
				params.getConnectionString(),
				Opts{},
				params.Files...,
			)
			if err != nil {
				log.Error(ctx, err, "unable to process migration up")
				return
			}
			if m != nil {
				if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
					log.Error(ctx, err, "unable to process migration up")
					return
				}
			}
		}
	}
}

func (p PreHookMigrationParams) getConnectionString() string {
	if p.PostgresDialect {
		return p.ConnectionString
	}
	return PrepareCockroach(p.ConnectionString)
}
