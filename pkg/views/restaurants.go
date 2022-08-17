package views

import (
	"context"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Get All Restaurants
*/
func GetRestaurants(c *gin.Context) {

	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

	var restaurants = []bson.M{}

	var response = MongoJsonResponse{}

	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(context.Background(), &restaurants); err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Type = "success"
	response.Data = restaurants
	response.Message = "Restaurants found"

	c.JSON(http.StatusOK, response)
}

func GetRestaurantByName(c *gin.Context) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

	var restaurants bson.M

	var restaurant []bson.M

	var response = MongoJsonResponse{}

	filter := bson.M{"name": "Shiro Lagos"}

	err := collection.FindOne(context.TODO(), filter).Decode(&restaurants)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant = append(restaurant, restaurants)

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	c.JSON(http.StatusOK, response)
}

func CreateRestaurant(c *gin.Context) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

	request := CreateRestaurants{}

	response := MongoJsonResponse{}

	err := c.BindJSON(&request)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.CreatedAt = config.GetCurrentTime()
	request.UpdatedAt = config.GetCurrentTime()

	filter := bson.M{
		"name":        request.Name,
		"description": request.Description,
		"address":     request.Address,
		"phone":       request.Phone,
		"email":       request.Email,
		"website":     request.Website,
		"image":       request.Image,
		"latitude":    request.Lat,
		"longitude":   request.Long,
		"rating":      request.Rating,
		"openingTime": request.OpeningTime,
		"price":       request.Price,
		"categories":  request.Categories,
		"tags":        request.Tags,
		"createdAt":   request.CreatedAt,
		"updatedAt":   request.UpdatedAt,
	}

	insertResult, err := collection.InsertOne(context.TODO(), filter)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Type = "success"
	response.Data = bson.M{"InsertID": insertResult.InsertedID.(primitive.ObjectID).Hex()}
	response.Message = "Restaurant created"

	c.JSON(http.StatusOK, response)

}
