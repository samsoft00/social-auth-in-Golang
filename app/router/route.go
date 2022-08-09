package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

type deps struct {
	fx.In

	GinEngine *gin.Engine
}

func SetupRoutes(d deps) {
	ginEngine := d.GinEngine

	// cors
	ginEngine.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Type", "Content-Length"},
	}))

	go func() {
		_ = http.ListenAndServe(":8080", ginEngine)
	}()
}
