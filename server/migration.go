package main

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigration() {
	databaseURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", databaseURL+"?sslmode=disable")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Printf("Failed to create migration driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/sql/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Printf("Failed to initialize migration: %v\n", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Printf("Migration failed: %v\n", err)
	}
}
