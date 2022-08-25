package views

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Get User by ID
*/
var usersCollection = database.OpenCollection(database.ConnectMongoDB(), "lagos_restaurants", "USERS")

func GetUserByID(c *gin.Context) {
	// client := database.ConnectMongoDB()

	// defer client.Disconnect(context.TODO())

	// collection := client.Database(restaurantDB).Collection("USERS")

	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var user bson.M
	request := GetUser{}
	var response = MongoJsonResponse{}

	if err := c.BindJSON(&request); err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Print("Request ID sent by client:", request.ID)

	id, err := primitive.ObjectIDFromHex(request.ID.Hex())
	config.CheckErr(err)

	config.Logs("info", "ID: "+id.Hex())

	filter := bson.M{"_id": id}
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Type = "success"
	response.Data = user
	response.Message = "User found"

	c.JSON(http.StatusOK, response)
}

func CreatNewUser(c *gin.Context) {
	// client := database.ConnectMongoDB()

	// defer client.Disconnect(context.TODO())

	// collection := client.Database(restaurantDB).Collection("USERS")
	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var user = CreatUser{}
	var response = MongoJsonResponse{}
	if err := c.BindJSON(&user); err != nil {
		config.Logs("error", err.Error())
		response.Type = "error"
		response.Message = "fullname, username, email, phone, password are required"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	config.Logs("info", "User: "+user.Fullname+" "+user.Username+" "+user.Email)
	user.Created_At = primitive.NewDateTimeFromTime(time.Now())
	user.Updated_At = primitive.NewDateTimeFromTime(time.Now())
	filter := bson.M{
		"fullname":   user.Fullname,
		"username":   user.Username,
		"avatar":     user.Avatar,
		"email":      user.Email,
		"password":   user.Password,
		"social":     user.Social,
		"role":       user.Role,
		"created_at": user.Created_At,
		"updated_at": user.Updated_At,
	}
	insertResult, err := usersCollection.InsertOne(context.TODO(), filter)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("insertResult: ", insertResult)
	response.Type = "success"
	response.Message = "User created"
	c.JSON(http.StatusOK, response)
}
