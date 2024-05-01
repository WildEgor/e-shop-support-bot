package pkg

import (
	"database/sql"
	"embed"
	"errors"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log/slog"
)

type Migrator struct {
	srcDriver source.Driver // Драйвер источника миграций.
}

func MustGetNewMigrator(sqlFiles embed.FS, dirName string) *Migrator {
	// Создаем новый драйвер источника миграций с встроенными SQL-файлами.
	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		panic(err)
	}

	return &Migrator{
		srcDriver: d,
	}
}

func (m *Migrator) ApplyMigrations(db *sql.DB) error {
	// Создаем экземпляр драйвера базы данных для PostgreSQL.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("unable to create db instance", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		return err
	}

	// Создаем новый экземпляр мигратора с использованием драйвера источника и драйвера базы данных PostgreSQL.
	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		slog.Error("unable to create migration", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		return err
	}

	// Закрываем мигратор в конце работы функции.
	defer func() {
		migrator.Close()
	}()

	// Применяем миграции.
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("unable to apply migrations", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		return err
	}

	return nil
}

const migrationsDir = "db/postgres/migrations"

//go:embed db/postgres/migrations/*.sql
var MigrationsFS embed.FS

func RunMigrate(connectionStr string) {

	migrator := MustGetNewMigrator(MigrationsFS, migrationsDir)

	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	err = migrator.ApplyMigrations(conn)
	if err != nil {
		panic(err)
	}

	slog.Debug("Migrations applied!!")
}
