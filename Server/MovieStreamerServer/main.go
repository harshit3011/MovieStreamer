package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	controller "github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/controllers"

)

func main() {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello My movie streamer!")
	})

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addMovie", controller.AddMovie())
	router.POST("/register",controller.RegisterUser())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
