package container

import (
	"context"

	"file-storage/internal/app"
	"file-storage/internal/config"

	"go.uber.org/fx"
)

func RegisterAppHooks(lc fx.Lifecycle, c *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return app.CheckFolders(c.Storage.Path)
		},
	})
}

var AppModule = fx.Module("app", fx.Options(
	fx.Invoke(RegisterAppHooks),
))
