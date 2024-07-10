package migration

import (
	"context"
	"embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/infiniteloopcloud/log"
	"github.com/spf13/afero"
)

type Opts struct {
	Path string
}

// Get will return the migration prepared for the certain files embedded into the binary
func Get(_ context.Context, conn string, opts Opts, embeds ...embed.FS) (*migrate.Migrate, error) {
	path := "files"
	if opts.Path != "" {
		path = opts.Path
	}
	afs, err := FromEmbeds(opts, embeds...)
	if err != nil {
		return nil, err
	}

	d, err := iofs.New(afero.NewIOFS(afs), path)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("iofs", d, conn)
}

func FromEmbeds(opts Opts, embeds ...embed.FS) (afero.Fs, error) {
	pathOuter := "files"
	if opts.Path != "" {
		pathOuter = opts.Path
	}
	afs := afero.NewMemMapFs()
	if err := afs.Mkdir(pathOuter, os.ModePerm); err != nil {
		return afs, err
	}
	for _, e := range embeds {
		err := fs.WalkDir(e, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				asd, err := e.Open(path)
				if err != nil {
					return err
				}
				c, err := io.ReadAll(asd)
				if err != nil {
					return err
				}
				err = afero.WriteFile(afs, path, c, os.ModePerm)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return afs, err
		}
	}
	return afs, nil
}

func MigrationTool(conn string, down bool, opts Opts, embeds ...embed.FS) {
	ctx := context.Background()
	m, err := Get(ctx, conn, opts, embeds...)
	if err != nil {
		log.Error(ctx, err, "unable to process migration")
		return
	}
	if m == nil {
		log.Error(ctx, errors.New("migration is nil"), "unable to process migration")
		return
	}
	if down {
		err = m.Down()
	} else {
		err = m.Up()
	}

	if err != nil {
		log.Error(ctx, err, "unable to process migration")
	}
	log.Info(ctx, "migration succeed")
}

func MigrationToolWithError(conn string, down bool, opts Opts, embeds ...embed.FS) error {
	ctx := context.Background()
	m, err := Get(ctx, conn, opts, embeds...)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("migration is nil")
	}
	if down {
		err = m.Down()
	} else {
		err = m.Up()
	}
	return err
}

func PrepareCockroach(conn string) string {
	return strings.Replace(conn, "postgres://", "cockroach://", 1)
}
