package container

import (
	"context"
	"fmt"

	"file-storage/internal/config"
	"file-storage/internal/logger"
	"file-storage/internal/postgresql"
	"file-storage/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func RegisterPostgreSQLHooks(lc fx.Lifecycle, c *config.Config, l *logger.Logger, pg *postgresql.PostgreSQL) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("trying to connect to postgresql")

			pool, err := pgxpool.New(ctx, c.PostgreSQL.Url)
			if err != nil {
				return fmt.Errorf("postgresql connection error: %w", err)
			}

			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("postgresql connection error: %w", err)
			}

			pg.Pool = pool

			l.Info("postgresql connection has been established")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			pg.Pool.Close()
			return nil
		},
	})
}

var PostgreSQLModule = fx.Module("pg", fx.Options(
	fx.Provide(postgresql.NewPostgreSQL),
	fx.Provide(fx.Annotate(
		// TODO: is it even possible to combine PostgreSQL + PostgreSQLReadiness and provide PostgreSQL with IsReady method?
		postgresql.NewPostgreSQLReadiness,
		fx.As(new(services.Readiness)),
		fx.ResultTags(`group:"readiness"`),
	)),
	fx.Invoke(RegisterPostgreSQLHooks),
))
