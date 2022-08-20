package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samsoft00/social-authentication-in-Golang/connector/tiktok"
	"go.uber.org/fx"
	"net/http"
)

type deps struct {
	fx.In

	TiktokController *tiktok.Controller
	GinEngine        *gin.Engine
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

	{
		t := d.GinEngine.Group("tiktok")

		t.GET("auth", d.TiktokController.Init)
		t.GET("callback", d.TiktokController.Callback)
	}

	go func() {
		_ = http.ListenAndServe(":8080", ginEngine)
	}()
}
