package repository

import "database/sql"

type BaseRepository[T any] struct {
	DB    *sql.DB
	Query T
}

func NewBaseRepository[T any](db *sql.DB, client T) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB:    db,
		Query: client,
	}
}
