package infrastructure

import (
	"context"
	"database/sql"
	"os"

	"github.com/XSAM/otelsql"
	_ "github.com/mattn/go-sqlite3"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

func SqlLite3DBConnect(ctx context.Context) *sql.DB {
	sqlite3Path := os.Getenv("SQLITE3_PATH")

	db, err := otelsql.Open("sqlite3", sqlite3Path,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithSpanOptions(otelsql.SpanOptions{DisableQuery: false}),
	)
	if err != nil {
		zap.L().Fatal("failed to connect to sqlite3", zap.Error(err))
	}

	_, err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemSqlite,
	))
	if err != nil {
		zap.L().Fatal("failed to register db stats metrics", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		zap.L().Fatal("failed to ping sqlite3", zap.Error(err))
	}

	zap.L().Info("success connected to sqlite3 in : " + sqlite3Path)
	return db
}
