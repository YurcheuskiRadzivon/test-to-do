package migrations

import (
	"database/sql"
	"fmt"

	migrator "github.com/rubenv/sql-migrate"
)

const (
	nameDBPG    = "postgres"
	defaultDirM = "sql/migrations"
)

type Migrator interface {
	Migrate() error
}

type PGMigrator struct {
	URL string
	Dir string
}

func NewPGMigrator(url string, dir string) *PGMigrator {
	return &PGMigrator{
		URL: url,
		Dir: dir,
	}
}

func (pgm *PGMigrator) Migrate() error {
	db, err := sql.Open(nameDBPG, pgm.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	migrations := &migrator.FileMigrationSource{
		Dir: pgm.Dir,
	}

	n, err := migrator.Exec(db, nameDBPG, migrations, migrator.Up)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
