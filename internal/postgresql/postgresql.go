package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PostgreSQL struct {
		Pool *pgxpool.Pool
	}

	PostgreSQLReadiness struct {
		pg *PostgreSQL
	}
)

func NewPostgreSQL() *PostgreSQL {
	return &PostgreSQL{}
}

func (r *PostgreSQLReadiness) IsReady(ctx context.Context) error {
	return r.pg.Pool.Ping(ctx)
}

func NewPostgreSQLReadiness(pg *PostgreSQL) *PostgreSQLReadiness {
	return &PostgreSQLReadiness{
		pg: pg,
	}
}
