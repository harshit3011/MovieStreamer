package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/controllers"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/middleware"
)

func SetupProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleWare())

	protected.GET("/movie/:imdb_id", controller.GetMovie())
	protected.POST("/addMovie", controller.AddMovie())
	protected.GET("/recommended", controller.GetRecommendedMovies())
}
