package controllers

import (
	"context"
	"fmt"
	"go-backend-ailanglearn/configs"
	"go-backend-ailanglearn/helpers"
	"go-backend-ailanglearn/models"
	"go-backend-ailanglearn/responses"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var messagesCollection *mongo.Collection = configs.GetColletion(configs.DB, "messages")

func PostMessage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//if request method is not POST
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed",
			},
		})
	}

	var message models.Message

	//validate the request body
	if err := c.BodyParser(&message); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//validate the request body
	validate := validator.New()

	if err := validate.Struct(message); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	var reqToken string = c.Get("Authorization")

	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) != 2 {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{
				"result": "Unauthorized",
			},
		})
	}

	reqToken = splitToken[1]

	//validate the token
	claims, status := helpers.ValidateToken(reqToken)

	if !status {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{
				"result": "Unauthorized",
			},
		})
	}

	message.Id = primitive.NewObjectID()
	message.Uid = claims.Uid
	fmt.Println("Message: " + message.Message)





	//get the user
	var user models.User





	//print userObjId


	userFilter := bson.M{"_id": claims.Uid}

	
	err := userCollection.FindOne(ctx, userFilter).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//insert the message in the messages field of the user
	user.Messages = append(user.Messages, message.Id)

	//update the user
	update := bson.M{
		"$set": bson.M{
			"messages": user.Messages,
		},
	}

	result, err := userCollection.UpdateOne(ctx, userFilter, update)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": "Error updating the user",
			},
		})
	}

	//insert the message in the messages collection
	_, err = messagesCollection.InsertOne(ctx, message)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"result": "Message sent",
		},
	})


}