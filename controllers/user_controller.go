package controllers

import (
	"context"
	"errors"
	"fmt"
	"go-backend-ailanglearn/configs"
	"go-backend-ailanglearn/helpers"
	"go-backend-ailanglearn/models"
	"go-backend-ailanglearn/responses"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = configs.GetColletion(configs.DB, "users")
var tokenCollection *mongo.Collection = configs.GetColletion(configs.DB, "tokens")
var noteCollection *mongo.Collection = configs.GetColletion(configs.DB, "notes")
var validate = validator.New()

func CreateUser(user *models.User, code int, token string, created time.Time) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//print all the user input
	fmt.Println("Id: " + user.Id.Hex())
	fmt.Println("Email: " + user.Email)
	fmt.Println("Password: " + user.Password)
	fmt.Println("Created_at: " + user.Created_at.String())

	//validate the user input
	if validationErr := validate.Struct(user); validationErr != nil {
		return bson.M{}, validationErr
	}

	newUser := models.User{
		Id:            primitive.NewObjectID(),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Title:         user.Title,
		Email:         user.Email,
		Password:      user.Password,
		Token:         user.Token,
		Refresh_token: user.Refresh_token,
		Created_at:    user.Created_at,
		Updated_at:    user.Updated_at,
		User_id:       user.User_id,
	}

	//check if the email already exists
	filter := bson.M{"email": user.Email}

	var result models.User

	if err := userCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No user found, good")
			if _, err := userCollection.InsertOne(ctx, newUser); err != nil {
				fmt.Println("error inserting user")
				return bson.M{}, err
			} else {

				if resultToken, err := HandleInsertToken(token, code, created); err != nil {
					fmt.Println("error inserting token")
					return bson.M{}, err
				} else {
					fmt.Println("Token inserted")
					return resultToken, nil
				}

			}

		} else {
			fmt.Println("error email already exists")
			return bson.M{}, err
		}

	} else {
		fmt.Println("email already exists")
		err := errors.New("email already exists")
		return bson.M{}, err
	}

}

func HandleInsertToken(token string, code int, created time.Time) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if code != 0 {
		data := bson.M{"token": token, "code": code, "created_at": created}

		_, errInsert := tokenCollection.InsertOne(ctx, data)

		if errInsert != nil {
			return bson.M{}, errInsert
		} else {
			return data, nil
		}
	} else if code == 0 {
		data := bson.M{"token": token, "created_at": created}

		_, errInsert := tokenCollection.InsertOne(ctx, data)

		if errInsert != nil {
			return bson.M{}, errInsert
		}
	}

	return bson.M{}, nil

}

func HandleSignup(c *fiber.Ctx) error {
	//if request method is not POST
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed`, bitchass bullushit",
			},
		})
	}

	var user models.User
	var result models.Response

	err := c.BodyParser(&user)
	fmt.Println("test1")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	/*
		userAgent := c.Get("User-Agent")
		//parse the user agent
		DeviceInfo := uasurfer.Parse(userAgent)
		fmt.Println("test2")

		var agent models.UserAgent

		agent.UserAgent = userAgent
		agent.OS.Name = DeviceInfo.OS.Name.String()
		agent.OS.Platform = DeviceInfo.OS.Platform.String()
		agent.OS.Version = string(rune(DeviceInfo.OS.Version.Major)) + "." + string(rune(DeviceInfo.OS.Version.Minor)) + "." + string(rune(DeviceInfo.OS.Version.Patch))
		agent.Browser.Name = DeviceInfo.Browser.Name.String()
		agent.Browser.Version = string(rune(DeviceInfo.Browser.Version.Major)) + "." + string(rune(DeviceInfo.Browser.Version.Minor)) + "." + string(rune(DeviceInfo.Browser.Version.Patch))
		agent.DeviceType = DeviceInfo.DeviceType.String()
	*/

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Id = primitive.NewObjectID()
	user.User_id = user.Id.Hex()
	token, refreshToken, _ := helpers.JWTTokenGenerator(user.Email, user.FirstName, user.LastName, user.User_id)
	user.Token = token
	user.Refresh_token = refreshToken
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(encryptedPassword)
	result.Token = token
	result.Expires_in = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	generatedCode := helpers.HandleCodeGenerator(6)
	code, _ := strconv.Atoi(generatedCode)

	//create the user and check for errors
	if creationResult, err := CreateUser(&user, code, result.Token, user.Created_at); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	} else {
		//send the user a verification code
		//helpers.SendVerificationCode(user.Email, generatedCode)

		return c.Status(http.StatusCreated).JSON(responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: &fiber.Map{
				"result": creationResult,
			},
		})
	}
}

func HandleSignin(c *fiber.Ctx) error {
	//if request method is not POST
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed`, bitchass bullushit",
			},
		})
	}

	var user models.AuthenticationModel
	var result models.Response

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	auth, email, fname, lname, userid := HandleAuthentication(user.Email, user.Password)
	token, _, _ := helpers.JWTTokenGenerator(email, fname, lname, userid)
	result.Token = token
	result.Expires_in = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	if !auth {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{
				"result": "Invalid credentials",
			},
		})
	} else {
		return c.Status(http.StatusOK).JSON(responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: &fiber.Map{
				"result": result,
			},
		})
	}

}

func HandleAuthentication(email string, password string) (bool, string, string, string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.AuthenticationModel

	errFind := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	decryptPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if decryptPassword != nil {
		return false, "", "", "", ""
	}

	if errFind != nil {
		return false, "", "", "", ""
	}

	return true, user.Email, user.FirstName, user.LastName, user.User_id
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
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
			"result": user,
		},
	})
}

func CreateNote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//if request method is not POST
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed`, bitchass bullushit",
			},
		})
	}

	var note models.Note

	err := c.BodyParser(&note)
	//get authorization header
	var reqToken string = c.Get("Authorization")
	fmt.Println(reqToken)
	splitToken := strings.Split(reqToken, "Bearer ")

	reqToken = splitToken[1]

	//validate the token
	claims, status := helpers.ValidateToken(reqToken)

	if !status {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{
				"result": "Invalid token",
			},
		})
	} else {
		note.User_id = claims.Uid

	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//print all the user input
	fmt.Println("Id: " + note.Id.Hex())
	fmt.Println("Title: " + note.Title)
	fmt.Println("Message: " + note.Message)
	fmt.Println("User_id: " + note.User_id)

	//validate the user input
	if validationErr := validate.Struct(note); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": validationErr.Error(),
			},
		})
	}

	note.Id = primitive.NewObjectID()
	note.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	note.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	result, err := noteCollection.InsertOne(ctx, note)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data: &fiber.Map{
			"result": result,
		},
	})

}

func GetProfile(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//if request method is not POST
	/*
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(responses.UserResponse{
			Status:  http.StatusMethodNotAllowed,
			Message: "error",
			Data: &fiber.Map{
				"result": "Method not allowed`, bitchass bullushit",
			},
		})
	}
*/

	//get authorization header
	var reqToken string = c.Get("Authorization")
	fmt.Println(reqToken)
	splitToken := strings.Split(reqToken, "Bearer ")

	reqToken = splitToken[1]

	//validate the token
	claims, status := helpers.ValidateToken(reqToken)

	if !status {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{
				"result": "Invalid token",
			},
		})
	} else {
		return c.Status(http.StatusOK).JSON(responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: &fiber.Map{
				"result": claims,
			},
		})

	}

}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10^time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{
				"result": validationErr.Error(),
			},
		})
	}

	update := bson.M{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Title":     user.Title,
		"title":     user.Title,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{
		"_id": objId,
	},
		bson.M{
			"$set": update,
		})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: &fiber.Map{
					"result": err.Error(),
				},
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"result": updatedUser,
		},
	})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data: &fiber.Map{
					"result": "user not found",
				},
			})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"result": "user deleted successfully",
		},
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	//find first 10 users
	results, err := userCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": -1}).SetLimit(10))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{
				"result": err.Error(),
			},
		})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: &fiber.Map{
					"result": err.Error(),
				},
			})
		}
		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"result": users,
		},
	})

}
