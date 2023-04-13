package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	Id	   	 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title,omitempty" validate:"required,min=2,max=100"`
	Message     string             `json:"message,omitempty" validate:"required,min=2,max=256"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}