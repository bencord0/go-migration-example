package models

import (
	"github.com/Masterminds/squirrel"
)

type SchemaMigrations struct {
	sql squirrel.StatementBuilderType
}

func SchemaMigrationsModel(db *Database) *SchemaMigrations {
	return &SchemaMigrations{
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db.db),
	}
}

type SchemaMigration struct {
	Version uint `json:"version"`
	Dirty bool `json:"dirty"`
}

func (sms *SchemaMigrations) CurrentVersion() (SchemaMigration, error) {
	q := sms.sql.
		Select("version", "dirty").
		From("schema_migrations").
		OrderBy("version ASC")

	rows, err := q.Query()
	if err != nil {
		return SchemaMigration{}, err
	}
	defer rows.Close()

	var migration SchemaMigration
	for rows.Next() {
		err = rows.Scan(
			&migration.Version,
			&migration.Dirty,
		)

		if err != nil {
			return SchemaMigration{}, err
		}
	}

	return migration, nil
}
