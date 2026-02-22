package infrastructure

import (
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func RunMigrations(db *sql.DB, migrations fs.FS) {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		zap.L().Fatal("failed to set goose dialect", zap.Error(err))
	}

	if err := goose.Up(db, "sql/migrations"); err != nil {
		zap.L().Fatal("failed to run migrations", zap.Error(err))
	}

	zap.L().Info("migrations applied successfully")
}
