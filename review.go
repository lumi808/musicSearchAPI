package main

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func reviewAddHandler(c *gin.Context) {
	userID, err := ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Convert userID to string
	userIDString := userID

	var payload ReviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not convert to ObjectID",
		})
		return
	}

	review := Review{
		AuthorID: objectID,
		Rating:  payload.Rating,
		Opinion: payload.Opinion,
	}

	// TODO: Check if album exists
	// TODO: Check if user has already reviewed this album

	updateAlbum := bson.D{{"$push", bson.D{{"reviews", review}}}}

	_, err = mongoClient.Database("music_search").Collection("albums").UpdateOne(
		context.TODO(), 
		bson.D{{"_id", payload.AlbumID}},
		updateAlbum,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert review to album",
			"error":   err,
		})
		return
	}

	var album Album
	err = mongoClient.Database("music_search").Collection("albums").
		FindOne(context.TODO(), bson.D{{"_id", payload.AlbumID}}).Decode(&album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not find album",
		})
		return
	}

	albumsCount := len(album.Reviews)
	totalRating := 0.0
	for _, review := range album.Reviews {
		totalRating += review.Rating
	}
	
	averageRating := totalRating / float64(albumsCount)

	updateAlbumRating := bson.D{{"$set", bson.D{{"rating", averageRating}}}}

	_, err = mongoClient.Database("music_search").Collection("albums").UpdateOne(
		context.TODO(),
		bson.D{{"_id", payload.AlbumID}},
		updateAlbumRating,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update album rating",
			"error":   err,
		})
		return
	}

	updateUser := bson.D{{"$push", bson.D{{"reviews", payload}}}}

	_, err = mongoClient.Database("music_search").Collection("users").UpdateOne(
		context.TODO(), 
		bson.D{{"_id", objectID}}, 
		updateUser,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert review to user",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review added successfully"})
}