package main

import (
	"backend/core/config"
	"backend/core/constants"
	"backend/core/pkg/errorsx"
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Running seed migrations...")

	cfg := config.Load(constants.ServiceScheduler)
	if err := runSeedMigrations(cfg); err != nil {
		log.Fatal(errorsx.WrapJSON(err, "Failed to run seed migrations"))
	}

	log.Println("Done.")
}

func runSeedMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", postgresDsn(cfg))
	if err != nil {
		return errorsx.Wrap(err, "Failed to open postgres connection")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_seeds",
	})
	if err != nil {
		return errorsx.Wrap(err, "Failed to create postgres driver")
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/seeds", "postgres", driver)
	if err != nil {
		return errorsx.Wrap(err, "Failed to create migrate instance")
	}

	if err := m.Down(); err != nil {
		log.Println(errorsx.WrapJSON(err, "Warning: down failed"))
	}

	_, _ = db.Exec(`TRUNCATE TABLE schema_seeds RESTART IDENTITY`)

	return m.Up()
}

func postgresDsn(cfg *config.Config) string {
	return "postgres://" +
		cfg.PGUser + ":" +
		cfg.PGPassword + "@" +
		cfg.PGHost + ":" +
		cfg.PGPort + "/" +
		cfg.PGDatabase +
		"?sslmode=disable"
}
