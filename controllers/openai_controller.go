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

	var responseMessage models.Message






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



	

	completions, completionsErr := helpers.CreateCompletionMessage(message.Message)

	if completionsErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": completionsErr.Error(),
			},
		})
	}

	responseMessage.Id = primitive.NewObjectID()
	responseMessage.Message = completions.Choices[0].Message.Content
	responseMessage.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	responseMessage.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	responseMessage.Right = false
	responseMessage.Uid = claims.Uid

	//Response should be like this
	/*
		{
		  "id": "chatcmpl-xxx",
		  "object": "chat.completion",
		  "created": 1678667132,
		  "model": "gpt-3.5-turbo-0301",
		  "usage": {
			"prompt_tokens": 13,
			"completion_tokens": 7,
			"total_tokens": 20
		  },
		  "choices": [
			{
			  "message": {
				"role": "assistant",
				"content": "\n\nThis is a test!"
			  },
			  "finish_reason": "stop",
			  "index": 0
			}
		  ]
		}
	*/

	message.Id = primitive.NewObjectID()
	message.Uid = claims.Uid
	message.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	message.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	
	fmt.Println("Message: " + message.Message)
	fmt.Println("Response: " + responseMessage.Message)





	//get the user
	var user models.User


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

	//insert the messages in the messages field of the user
	
	user.Messages = append(user.Messages, message.Id, responseMessage.Id)

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

	//create a json with a "message": message and "response": "Message posted successfully"


	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"result": fiber.Map{
				"message": message,
				"response": responseMessage,
			},
		},
	})


}


func GetMessages(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//if request method is not GET
	if c.Method() != "GET" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed",
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

	//get the user
	var user models.User

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

	if len(user.Messages) == 0 {
		return c.Status(http.StatusOK).JSON(responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: &fiber.Map{
				"result": "No messages",
			},
		})
	}

	//get the messages
	var messages []models.Message

	messagesFilter := bson.M{"_id": bson.M{"$in": user.Messages}}

	cursor, err := messagesCollection.Find(ctx, messagesFilter)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	if err = cursor.All(ctx, &messages); err != nil {
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
			"result": messages,
		},
	})
}
