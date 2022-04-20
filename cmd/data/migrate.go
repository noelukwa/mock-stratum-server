package data

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed migrations
var migrations embed.FS

func RunMigration(dsn string) error {
	log.Println("Running migration")
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, dsn)

	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("Done with  migration")
	}()
	return m.Up()
}
