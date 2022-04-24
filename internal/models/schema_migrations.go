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
	Version int `json:"version"`
	Dirty bool `json:"dirty"`
}

func (sms *SchemaMigrations) List() ([]SchemaMigration, error) {
	q := sms.sql.
		Select("version", "dirty").
		From("schema_migrations").
		OrderBy("version ASC")

	rows, err := q.Query()
	if err != nil {
		return []SchemaMigration{}, err
	}
	defer rows.Close()

	var migrations []SchemaMigration
	for rows.Next() {
		var migration SchemaMigration
		err = rows.Scan(
			&migration.Version,
			&migration.Dirty,
		)

		if err != nil {
			return []SchemaMigration{}, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}
