package main

import (
	"errors"
	"flag"
	"fmt"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций Postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Драйвер для получения миграция из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, dbUsername, dbPassword, host, port, dbName, query string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&dbUsername, "db-username", "", "db username")
	flag.StringVar(&dbPassword, "db-password", "", "db user password")
	flag.StringVar(&host, "db-host", "localhost", "db host")
	flag.StringVar(&port, "db-port", "5432", "db port")
	flag.StringVar(&dbName, "db-name", "", "db name")
	flag.StringVar(&query, "db-query", "", "query to db")
	flag.Parse()

	if dbName == "" {
		panic("db-name is required")
	}

	if migrationsPath == "" {
		panic("migrations_path is required")
	}

	if dbUsername == "" {
		panic("username is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgresql://%s@%s:%s/%s%s",
			dbUsername,
			host,
			port,
			dbName,
			query),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
