package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
}

type SQLStore struct {
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
	}
}
