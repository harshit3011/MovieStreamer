package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/controllers"
)

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controller.GetMovies())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
	router.POST("/logout", controller.LogoutHandler())
}
