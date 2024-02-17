package container

import (
	"context"

	"file-storage/internal/config"
	"file-storage/internal/logger"
	"file-storage/internal/postgresql"
	"file-storage/internal/services"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
)

func RegisterPostgreSQLHooks(lc fx.Lifecycle, c *config.Config, l *logger.Logger, pg *postgresql.PostgreSQL) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			return pg.Start(ctx, &postgresql.PostgreSQLStartOptions{
				ConnectionUrl: c.PostgreSQL.Url,
				BeforeConnect: func(ctx context.Context, cc *pgx.ConnConfig) error {
					l.Info("trying to connect to postgresql...")
					return nil
				},
				AfterConnect: func(ctx context.Context, c *pgx.Conn) error {
					l.Info("successfully connected to postgresql")
					return nil
				},
			})
		},
		OnStop: func(ctx context.Context) error {
			return pg.Stop(ctx)
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
