package main

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	Name string 				`json:"name,omitempty" bson:"name,omitempty"`
	Surname string 				`json:"surname,omitempty" bson:"surname,omitempty"`
	Avatar   string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Reviews  []Review           `json:"reviews,omitempty" bson:"reviews,omitempty"`
}

type LoginPayload struct {
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
}

type Review struct {
	AuthorID primitive.ObjectID `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Rating  float64    `json:"rating,omitempty" bson:"rating,omitempty"`
	Opinion string `json:"opinion,omitempty" bson:"opinion,omitempty"`
}

type ReviewPayload struct {
	AlbumID primitive.ObjectID `json:"album_id,omitempty" bson:"album_id,omitempty"`
	Rating  float64    `json:"rating,omitempty" bson:"rating,omitempty"`
	Opinion string `json:"opinion,omitempty" bson:"opinion,omitempty"`
}

type Track struct {
	Title string  `json:"title,omitempty" bson:"title,omitempty"`
	Length string `json:"length,omitempty" bson:"length,omitempty"`
}

type Album struct {
	ID primitive.ObjectID           `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string 		            `json:"title,omitempty" bson:"title,omitempty"`
	Artists []string                `json:"artists,omitempty" bson:"artists,omitempty"`
	DateReleased primitive.DateTime `json:"date_released,omitempty" bson:"date_released,omitempty"`
	Description string              `json:"description,omitempty" bson:"description,omitempty"`
	Cover string                    `json:"cover,omitempty" bson:"cover,omitempty"`
	Tracklist []Track               `json:"tracklist,omitempty" bson:"tracklist,omitempty"`
	Reviews []Review                `json:"reviews,omitempty" bson:"reviews,omitempty"`
	AverageRating float64 			`json:"rating,omitempty" bson:"rating,omitempty"`
}

type AlbumPayload struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string          `json:"title,omitempty" bson:"title,omitempty"`
	Artists []string      `json:"artists,omitempty" bson:"artists,omitempty"`
	DateReleased string   `json:"date_released,omitempty" bson:"date_released,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Cover string          `json:"cover,omitempty" bson:"cover,omitempty"`
	Tracklist []Track     `json:"tracklist,omitempty" bson:"tracklist,omitempty"`
	Reviews []Review      `json:"reviews,omitempty" bson:"reviews,omitempty"`
	AverageRating float64 `json:"average_rating,omitempty" bson:"average_rating,omitempty"`
}

type Session struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AccessToken string `json:"access_token,omitempty" bson:"access_token,omitempty"`
}