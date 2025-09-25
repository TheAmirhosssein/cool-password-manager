package database

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/TheAmirhosssein/cool-password-manage/pkg/testdocker"
	"github.com/TheAmirhosssein/goose/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var dbPool *pgxpool.Pool

func initDB(ctx context.Context, dsn string) error {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	dbPool = pool
	return nil
}

func GetDb(ctx context.Context) *pgxpool.Pool {
	if dbPool == nil {
		conf := config.GetConfig()
		err := initDB(ctx, conf.DB.URL)
		if err != nil {
			panic(err.Error())
		}
	}
	return dbPool
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func migrate(db *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.UpPGX(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func SetupTestDB(ctx context.Context) (string, *pgxpool.Pool) {
	// Start transaction
	name, err := testdocker.GenerateName()
	if err != nil {
		panic("error getting name")
	}

	port, err := testdocker.StartPostgresContainer(ctx, name, name)
	if err != nil {
		panic("error starting database")
	}

	db, err := testdocker.GetTestDB(ctx, port)
	if err != nil {
		panic("error getting db")
	}

	err = migrate(db)
	if err != nil {
		panic("error getting db")
	}

	return name, db
}
