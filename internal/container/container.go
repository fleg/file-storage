package container

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"file-storage/internal/logger"
	"file-storage/internal/router"
	"file-storage/internal/services"
)

func registerLogger(l *logger.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{
		Logger: l.Desugar(),
	}
}

func Run() {
	appModule := fx.Options(
		fx.WithLogger(registerLogger),

		ConfigModule,
		LoggerModule,
		PostgreSQLModule,
		ServerModule,

		fx.Provide(fx.Annotate(services.NewHealthService, fx.ParamTags(`group:"readiness"`))),

		fx.Invoke(router.Bind),
	)

	app := fx.New(appModule)

	app.Run()
}
