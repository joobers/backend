package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client int

const (
	Android Client = iota
	iOS
	Web
	Chromium
	Firefox
)

type Message struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Body       string             `json:"body,omitempty" bson:"body,omitempty"`
	ClientType Client             `json:"client,omitempty" bson:"client,omitempty"`
	Tags       []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

// jwt.StandardClaims is an embedded type to provide fields like expiration time
type JwtClaims struct {
	ID primitive.ObjectID `json:"_id"`
	jwt.StandardClaims
}
