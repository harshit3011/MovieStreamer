package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/routes"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello My movie streamer!")
	})

	routes.SetupProtectedRoutes(router)
	routes.SetupUnProtectedRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
