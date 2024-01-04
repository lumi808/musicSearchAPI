package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func albumAddHandler(c *gin.Context){
	var payload AlbumPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	releaseDate, err := time.Parse("RFC3339", payload.DateReleased)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format",
		})
	}

	tracklist := make([]Track, len(payload.Tracklist))

	for i, track := range payload.Tracklist {
		tracklist[i] = Track{
			Title: track.Title,
			Length: track.Length,
		}
	}

	album := Album{
		Title: payload.Title,
		Artists: payload.Artists,
		DateReleased: primitive.NewDateTimeFromTime(releaseDate),
		Description: payload.Description,
		Cover: payload.Cover,
		Tracklist: tracklist,
		Reviews: payload.Reviews,
		AverageRating: payload.AverageRating,
	}

	//insert album into database
	result, err := mongoClient.Database("music_search").Collection("albums").InsertOne(context.TODO(), album)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not insert album",
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Album added successfully",
		"album": result,
	})
}

func albumGetAllHandler(c *gin.Context){
	var albums []Album

	cursor, err := mongoClient.Database("music_search").Collection("albums").Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var album Album
		err := cursor.Decode(&album)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, album)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Albums retrieved successfully",
		"albums": albums,
	})
}

func albumGetHandler(c *gin.Context){
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	var album Album

	err = mongoClient.Database("music_search").Collection("albums").
		FindOne(context.TODO(), bson.D{{"_id", objectID}}).Decode(&album)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not find album",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Album retrieved successfully",
		"album": album,
	})
}

func albumUpdateHandler(c *gin.Context){
	var payload AlbumPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	update := bson.D{}
	var updatedFields string = ""

	if payload.Title != "" {
		update = append(update, bson.E{"title", payload.Title})
		updatedFields += "title "
	}

	if payload.Artists != nil {
		update = append(update, bson.E{"artists", payload.Artists})
		updatedFields += "artists "
	}

	if payload.DateReleased != "" {
		releaseDate, err := time.Parse("RFC3339", payload.DateReleased)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid date format",
			})
		}

		update = append(update, bson.E{"date_released", primitive.NewDateTimeFromTime(releaseDate)})
		updatedFields += "date_released "
	}

	if payload.Description != "" {
		update = append(update, bson.E{"description", payload.Description})
		updatedFields += "description "
	
	}

	if payload.Cover != "" {
		update = append(update, bson.E{"cover", payload.Cover})
		updatedFields += "cover "
	}

	if payload.Tracklist != nil {
		tracklist := make([]Track, len(payload.Tracklist))

		for i, track := range payload.Tracklist {
			tracklist[i] = Track{
				Title: track.Title,
				Length: track.Length,
			}
		}

		update = append(update, bson.E{"tracklist", tracklist})
		updatedFields += "tracklist "
	}

	albumID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid album ID",
		})
		return
	}

	_, err = mongoClient.Database("music_search").Collection("albums").UpdateOne(
		context.TODO(), 
		bson.D{{"_id", objectID}}, 
		bson.D{{"$set", update}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update album",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Album updated successfully",
		"updated_fields": updatedFields,
	})
}