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
	return App(Config{
		ServiceName: "social-login",
		CusTomFxOptions: []fx.Option{
			// controller/service goes here <-
			fx.Provide(),
			fx.Invoke(router.SetupRoutes),
		},
	})
}

func App(config Config) []fx.Option {
	options := []fx.Option{
		fx.Provide(NewZapLogger),
		fx.Provide(SetupGin),
	}

	options = append(options, config.CusTomFxOptions...)
	return options
}

// SetupGin - initialize gin framework
func SetupGin(logger *zap.Logger) *gin.Engine {
	return gin.New()
}

// NewZapLogger -logger
func NewZapLogger() *zap.Logger {
	var logger *zap.Logger
	logger, _ = zap.NewDevelopment()

	return logger
}
