package admin

import (
	"net/http"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bencord0/go-migration-example/internal/models"
	. "github.com/bencord0/webframework"
)

func ListMigrations(db *models.Database) func(*http.Request) *http.Response {
	schemaMigrations := models.SchemaMigrationsModel(db)

	return func(req *http.Request) *http.Response {
		currentVersion, err := schemaMigrations.CurrentVersion()
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		migrationList := make(map[string]interface{})
		migrationList["version"] = currentVersion.Version
		migrationList["dirty"] = currentVersion.Dirty

		migrator, err := source.Open("file://migrations")
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		var migrations []map[string]interface{}
		for version, _ := migrator.First(); version != 0; version, _ = migrator.Next(version) {
			migration := make(map[string]interface{})
			migration["version"] = version
			if version < currentVersion.Version {
				_, ident, err := migrator.ReadDown(version)
				migration["identifier"] = ident
				if err != nil {
					continue
				}

				migration["status"] = "applied"
			} else {
				_, ident, err := migrator.ReadUp(version)
				migration["identifier"] = ident
				if err != nil {
					continue
				}

				migration["status"] = "pending"
				if version == currentVersion.Version {
					migration["status"] = "current"
				}
			}

			migrations = append(migrations, migration)
		}

		migrationList["migrations"] = migrations
		return JSONResponse(migrationList, http.StatusOK, nil)
	}
}

func ApplyMigrations(db *models.Database) func(*http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		if req.Method != "POST" {
			return TextResponse("method not allowd", http.StatusMethodNotAllowed, nil)
		}

		instance, err := db.GoMigrateInstance()
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres",
			instance,
		)
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		err = migrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		return TextResponse("migration successful", http.StatusOK, nil)
	}
}
