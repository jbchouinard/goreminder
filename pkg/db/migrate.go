package db

import (
	"context"
	"embed"
	"io/fs"
	"os"

	pgx4 "github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

const versionTable = "public.version"
const targetVersion = 1

type EmbeddedMigratorFS struct {
	fs *embed.FS
}

func (m *EmbeddedMigratorFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	entries, err := m.fs.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func (m *EmbeddedMigratorFS) ReadFile(filename string) ([]byte, error) {
	return m.fs.ReadFile(filename)
}

func (m *EmbeddedMigratorFS) Glob(pattern string) ([]string, error) {
	return fs.Glob(m.fs, pattern)
}

func NewMigrator(ctx context.Context, conn *pgx4.Conn, versionTable string) (*migrate.Migrator, error) {
	return migrate.NewMigratorEx(
		ctx, conn, versionTable,
		&migrate.MigratorOptions{
			DisableTx:  false,
			MigratorFS: &EmbeddedMigratorFS{&migrations},
		})
}

func Migrate(ctx context.Context, databaseURL string) error {
	conn, err := pgx4.Connect(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	migrator, err := NewMigrator(ctx, conn, versionTable)
	if err != nil {
		return err
	}
	migrator.LoadMigrations("migrations")
	log.Info().Msgf("migrating database to version %d", targetVersion)
	return migrator.MigrateTo(ctx, targetVersion)
}
