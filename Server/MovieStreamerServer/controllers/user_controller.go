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
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil

}

func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userCollection *mongo.Collection = databases.OpenCollection("users", databases.Client)

		count, err := userCollection.CountDocuments(c, bson.D{{Key: "email", Value: user.Email}})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}
		if count > 0 {
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(c, user)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		ctx.JSON(http.StatusCreated, result)

	}

}