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

	cfg := config.Load(constants.ServiceServer)
	if err := runSeedMigrations(cfg); err != nil {
		log.Fatalf("%s", errorsx.JSONTrace(errorsx.Errorf("Failed to run seed migrations: %v", err)))
	}

	log.Println("Done.")
}

func runSeedMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", postgresDsn(cfg))
	if err != nil {
		return errorsx.Error(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_seeds",
	})
	if err != nil {
		return errorsx.Error(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/seeds", "postgres", driver)
	if err != nil {
		return errorsx.Error(err)
	}

	if err := m.Down(); err != nil {
		log.Printf("%s", errorsx.JSONTrace(errorsx.Errorf("Warning: down failed: %v", err)))
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
