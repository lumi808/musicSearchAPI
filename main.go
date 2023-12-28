package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var JWTSecret = []byte("your-secret-key")
var uri = "mongodb+srv://kairat:123@cluster0.umibhlh.mongodb.net/?retryWrites=true&w=majority"
var mongoClient *mongo.Client

func init(){
	if err:= connectToMongoDB(); err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	log.Println("Connected to MongoDB!")
}

func main() {
	fmt.Println("Hello, World!")
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hi, this is musicSearchAPI Backend Server",
		})
	})

	//account related endpoints
	r.POST("/register", accountRegisterHandler)
	r.POST("/login", accountLoginHandler)
	r.GET("/user/:id", accountGetUserHandler)
	
	r.Run() // listen and serve on 0.0.0.0:8080
}



func connectToMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
        return err
    }

	err = client.Ping(context.TODO(), nil)

	if err != nil {
        return err
    }

	mongoClient = client
	return nil
}