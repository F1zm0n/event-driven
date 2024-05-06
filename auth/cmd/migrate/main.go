package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	db := sqlx.MustOpen("postgres", os.Getenv("DSN"))
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{DatabaseName: "auth"})
	if err != nil {
		panic(err)
	}
	_, filename, _, _ := runtime.Caller(0)
	migrationPath := "file://" + filepath.Join(filepath.Dir(filename), "../../migrations")
	m, err := migrate.NewWithDatabaseInstance(migrationPath, "auth", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
}
