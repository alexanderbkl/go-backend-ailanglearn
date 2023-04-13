package helpers

import (
	"context"
	"fmt"
	"go-backend-ailanglearn/configs"
	"go-backend-ailanglearn/models"
	"log"
	"time"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY string = "SECRET_KEY"

func JWTTokenGenerator(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, status bool) {
	var updatedToken models.Response
	var userCollection *mongo.Collection = configs.GetColletion(configs.DB, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {

		return nil, false
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {

		return nil, false
	}

	//print claim expiration time in readable format
	fmt.Println("expires at")
	fmt.Println(time.Unix(claims.ExpiresAt, 0))

	if claims.ExpiresAt < time.Now().Local().Unix() {

		return nil, false
	}

	updatedToken.Token = signedToken

	updatedToken.Expires_in = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	//update token in db
	userCollection.UpdateOne(ctx, models.User{Token: signedToken}, updatedToken)

	return claims, true
}