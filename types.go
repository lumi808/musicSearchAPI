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
	Rating  int    `json:"rating,omitempty" bson:"rating,omitempty"`
	Opinion string `json:"opinion,omitempty" bson:"opinion,omitempty"`
}
