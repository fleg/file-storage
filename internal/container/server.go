package container

import (
	"context"
	"net/http"

	"file-storage/internal/config"
	"file-storage/internal/server"
	"file-storage/internal/logger"
	"go.uber.org/fx"
)

func RegisterServerHooks(lc fx.Lifecycle, c *config.Config, l *logger.Logger, s *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := s.Start(&server.StartOptions{Port: c.Server.Port}); err != nil && err != http.ErrServerClosed {
					l.Fatalf("can't start http server: %+v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})
}

var ServerModule = fx.Options(
	fx.Provide(server.New),
	fx.Invoke(RegisterServerHooks),
)
