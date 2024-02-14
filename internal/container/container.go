package container

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"go-files/internal/logger"
)

func registerLogger(l *logger.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{
		Logger: l.Desugar(),
	}
}

func Run() {
	appModule := fx.Options(
		ConfigModule,
		LoggerModule,

		fx.WithLogger(registerLogger),
		ServerModule,
	)

	app := fx.New(appModule)

	app.Run()
}
