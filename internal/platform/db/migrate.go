package db

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(databaseURL, migrationsDir string) error {
	src := migrationsDir
	if !strings.HasPrefix(src, "file://") {
		src = "file://" + strings.ReplaceAll(migrationsDir, "\\", "/")
	}
	dsn := databaseURL
	if u, err := url.Parse(databaseURL); err == nil {
		q := u.Query()
		if q.Get("sslmode") == "" {
			q.Set("sslmode", "disable")
			u.RawQuery = q.Encode()
			dsn = u.String()
		}
	}
	m, err := migrate.New(src, dsn)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
