package admin

import (
	"net/http"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bencord0/go-migration-example/internal/models"
	. "github.com/bencord0/webframework"
)

func ListMigrations(db *models.Database) func(*http.Request) *http.Response {
	schemaMigrations := models.SchemaMigrationsModel(db)

	return func(req *http.Request) *http.Response {
		allMigrations, err := schemaMigrations.List()
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		return JSONResponse(allMigrations, http.StatusOK, nil)
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
