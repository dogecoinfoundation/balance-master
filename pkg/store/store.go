package store

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB             *sql.DB
	backend        string
	migrationsPath string
}

func NewStore(dbUrl string, migrationsPath string) (*Store, error) {
	u, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "memory" {
		sqlite, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			return nil, err
		}

		return &Store{DB: sqlite, backend: "sqlite", migrationsPath: migrationsPath}, nil
	} else if u.Scheme == "sqlite" {
		var url string
		if u.Host == "" {
			url = u.Path
		} else {
			url = u.Host
		}

		sqlite, err := sql.Open("sqlite3", url)
		if err != nil {
			return nil, err
		}

		return &Store{DB: sqlite, backend: "sqlite", migrationsPath: migrationsPath}, nil
	} else if u.Scheme == "postgres" {
		postgres, err := sql.Open("postgres", dbUrl)
		if err != nil {
			return nil, err
		}
		return &Store{DB: postgres, backend: "postgres", migrationsPath: migrationsPath}, nil
	}

	return nil, fmt.Errorf("unsupported database scheme: %s", u.Scheme)
}

func (s *Store) Migrate() error {
	driver, err := s.getMigrationDriver()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+s.migrationsPath, s.backend, driver)
	if err != nil {
		return err
	}

	return m.Up()
}

func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory, cannot find go.mod
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func MigrationsPath() (string, error) {
	root, err := ProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "db", "migrations"), nil
}

func (s *Store) getMigrationDriver() (database.Driver, error) {
	if s.backend == "postgres" {
		driver, err := postgres.WithInstance(s.DB, &postgres.Config{})
		if err != nil {
			return nil, err
		}

		return driver, nil
	}

	if s.backend == "sqlite" {
		driver, err := sqlite.WithInstance(s.DB, &sqlite.Config{})
		if err != nil {
			return nil, err
		}

		return driver, nil
	}

	return nil, fmt.Errorf("unsupported database scheme: %s", s.backend)
}

func (s *Store) Close() error {
	fmt.Println("Closing store")
	return s.DB.Close()
}
