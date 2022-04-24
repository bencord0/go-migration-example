package models

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Database struct {
	url *url.URL
	db *sql.DB
}

func NewDatabase(databaseUrl *url.URL) (*Database, error) {
	c, err := pgx.ParseConfig(databaseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("parsing database url: %w", err)
	}

	database := stdlib.OpenDB(*c)
	return &Database{
		url: databaseUrl,
		db: database,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GoMigrateInstance() (database.Driver, error) {
	return postgres.WithInstance(d.db, new(postgres.Config))
}
