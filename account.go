package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func accountRegisterHandler(c *gin.Context){
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	if userExists, err := isUserExists(user.Username, user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error checking user existence",
			"error":   err,
		})
		return
	} else if userExists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User with the same username or email already exists",
		})
		return
	}

	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not hash password",
		})
		return
	}

	user.Password = hashedPassword

	result, err := mongoClient.Database("music_search").Collection("users").InsertOne(context.TODO(), user)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert user",
			"error": err,
		})
		return
	}

	token, err := generateJWTToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not generate JWT token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"token": token,
		"result": result,
	})
}

func accountLoginHandler(c *gin.Context){
	var user User
	var payload LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	err := mongoClient.Database("music_search").Collection("users").
		FindOne(context.TODO(), bson.D{{"email", payload.Email},}).Decode(&user)

	if err != nil {
		fmt.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not find user",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := generateJWTToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func accountGetUserHandler(c *gin.Context) {
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user User
	err = mongoClient.Database("music_search").Collection("users").
		FindOne(context.TODO(), bson.D{{"_id", objectID}}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func isUserExists(username, email string) (bool, error) {
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"username", username}},
			bson.D{{"email", email}},
		}},
	}

	count, err := mongoClient.Database("music_search").Collection("users").CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}