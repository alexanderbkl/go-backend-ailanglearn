package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id	   	 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName    string             `json:"first_name,omitempty" validate:"required,min=2,max=100"`
	LastName     string             `json:"last_name,omitempty" validate:"required,min=2,max=100"`
	Title	 string  			`json:"title,omitempty" bson:"title,omitempty"`
	Email         string             `json:"email" validate:"email,required"`
	Password string				`json:"password" bson:"password,omitempty" validate:"required,min=6,max=100"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

type Response struct {
	Token string `json:"token"`
	Expires_in int64 `json:"expires_in"`
}


type UserAgent struct {
	UserAgent string `json:"user_agent"`
	Browser   struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"browser"`
	OS struct {
		Platform string `json:"platform"`
		Name     string `json:"name"`
		Version  string `json:"version"`
	} `json:"os"`
	DeviceType string `json:"device_type"`
}

type AuthenticationModel struct {
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Email      string    `json:"email" validate:"required" bson:"email"`
	Password   string    `json:"password" validate:"required" bson:"password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    string    `json:"user_id"`
}

type Code struct {
	Code int `json:"code"`
}