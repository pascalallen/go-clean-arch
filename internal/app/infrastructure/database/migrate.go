package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/database/seeders"
)

func RunMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create database driver: %v", err)
	}

	migrationPath := "file://internal/app/infrastructure/database/migrations"

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run up migrations: %v", err)
	}

	fmt.Println("Migrations ran successfully")

	fmt.Println("Starting database seeding...")
	seeders.SeedDatabase(db)
	fmt.Println("Database seeding completed successfully.")
}
