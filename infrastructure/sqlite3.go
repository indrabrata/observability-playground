package infrastructure

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func SqlLite3Connect() *sql.DB {
	sqlite3Path := os.Getenv("SQLITE3_PATH")
	db, err := sql.Open("sqlite3", sqlite3Path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to sqlite3")
	}

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping sqlite3")
	}

	log.Info().Msg("success connected to sqlite3 in : " + sqlite3Path)
	return db
}
