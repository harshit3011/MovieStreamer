package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/databases"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/models"
	"github.com/harshit3011/MovieStreamer/Server/MovieStreamerServer/utils"
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

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalide input data"})
			return
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var userCollection *mongo.Collection = databases.OpenCollection("users", databases.Client)

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userLogin.Email}}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, databases.Client)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: token,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:    foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Role:      foundUser.Role,
			Token:     token,
			//RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})

	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var UserLogout struct {
			UserId string `json:"user_id"`
		}

		err := c.ShouldBindJSON(&UserLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request:", UserLogout.UserId)

		err = utils.UpdateAllTokens(UserLogout.UserId, "", "", databases.Client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}
