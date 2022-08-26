package routers

import (
	"FunNow/url-shortener/constants"
	"FunNow/url-shortener/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{
		constants.BaseUrl,
	}
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	r.POST("/shorten", controllers.PostShortenUrlHandler)

	return r
}
