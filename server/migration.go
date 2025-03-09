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
		fmt.Println("Failed to open database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/sql/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Println("Failed to initialize migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("Migration failed: %v", err)
	}
}
