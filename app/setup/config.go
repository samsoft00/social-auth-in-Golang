package setup

import (
	"github.com/gin-gonic/gin"
	"github.com/samsoft00/social-authentication-in-Golang/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	ServiceName     string
	CusTomFxOptions []fx.Option
}

func GetOptions() []fx.Option {
	return SetupApp(Config{
		ServiceName: "social-login",
		CusTomFxOptions: []fx.Option{
			// add your controller/service here <-
			fx.Provide(),
			fx.Invoke(router.SetupRoutes),
		},
	})
}

func SetupApp(config Config) []fx.Option {
	options := []fx.Option{
		fx.Provide(NewZapLogger),
		fx.Provide(func(logger *zap.Logger) *gin.Engine { return gin.New() }),
	}

	options = append(options, config.CusTomFxOptions...)
	return options
}

// NewZapLogger -logger
func NewZapLogger() *zap.Logger {
	var logger *zap.Logger
	logger, _ = zap.NewDevelopment()

	return logger
}
