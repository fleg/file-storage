package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PostgreSQL struct {
		Pool *pgxpool.Pool
	}

	PostgreSQLStartOptions struct {
		ConnectionUrl string
		BeforeConnect func(context.Context, *pgx.ConnConfig) error
		AfterConnect  func(context.Context, *pgx.Conn) error
	}

	PostgreSQLReadiness struct {
		pg *PostgreSQL
	}
)

func (pg *PostgreSQL) Start(ctx context.Context, options *PostgreSQLStartOptions) error {
	c, err := pgxpool.ParseConfig(options.ConnectionUrl)
	if err != nil {
		return fmt.Errorf("postgresql parse connection url error: %w", err)
	}

	c.BeforeConnect = options.BeforeConnect
	c.AfterConnect = options.AfterConnect

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return fmt.Errorf("postgresql connection error: %w", err)
	}

	pg.Pool = pool

	// trigger connection but ignore the result
	pool.Ping(ctx)

	return nil
}

func (pg *PostgreSQL) Stop(ctx context.Context) error {
	pg.Pool.Close()
	return nil
}

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
