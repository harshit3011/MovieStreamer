package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/databases"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

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

func GetMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx, 100*time.Second)
		defer cancel()

		movieID := ctx.Param("imdb_id")

		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movieCollection *mongo.Collection = databases.OpenCollection("movies", databases.Client)

		var movie models.Movie

		err := movieCollection.FindOne(c, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		ctx.JSON(http.StatusOK, movie)

	}
}

func AddMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := ctx.ShouldBindJSON(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		var movieCollection *mongo.Collection = databases.OpenCollection("movies", databases.Client)

		result, err := movieCollection.InsertOne(c, movie)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		ctx.JSON(http.StatusCreated, result)

	}
}