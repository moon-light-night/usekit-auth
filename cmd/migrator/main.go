package main

import (
	"errors"
	"flag"
	"fmt"
	// библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// драйвер для выполнения миграций SQLite3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "", "Path to the storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Path to a tables containing the migration files")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	migrator, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations found")
			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied successfully")
}
