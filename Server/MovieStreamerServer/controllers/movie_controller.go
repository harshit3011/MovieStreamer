package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/databases"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetMovies() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		c, cancel:= context.WithTimeout(ctx, 100*time.Second)
		defer cancel()
		var movieCollection *mongo.Collection = databases.OpenCollection("movies", databases.Client)

		cursor, err := movieCollection.Find(c, bson.D{})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(c)

		var movies []models.Movie

		if err = cursor.All(c, &movies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		ctx.JSON(http.StatusOK, movies)


	}
}