package views

import (
	"context"
	"net/http"
	"time"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Get All Restaurants
*/

var restaurantsCollection = database.OpenCollection(database.ConnectMongoDB(), DB, RESTAURANTS)

func GetRestaurants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var restaurants = []bson.M{}

	var response = MongoJsonResponse{}

	filter := bson.M{}

	opts := options.Find().SetProjection(bson.M{"title": 1, "description": 1, "address": 1, "images": 1, "rating": 1, "slug": 1})

	cursor, err := restaurantsCollection.Find(ctx, filter, opts)
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

func GetRestaurantByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var restaurant bson.M

	var request = GetRestaurantReviewsType{}

	var response = MongoJsonResponse{}

	if err := c.BindJSON(&request); err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(request.ID.Hex())

	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": id}

	err = restaurantsCollection.FindOne(ctx, filter).Decode(&restaurant)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	c.JSON(http.StatusOK, response)
}

func CreateRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

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
		"title":       request.Title,
		"description": request.Description,
		"slug":        request.Slug,
		"address":     request.Address,
		"phone":       request.Phone,
		"email":       request.Email,
		"website":     request.Website,
		"images":      request.Images,
		"rating":      request.Rating,
		"openingTime": request.OpeningTime,
		"currency":    request.Currency,
		"avgPrice":    request.Price,
		"categories":  request.Categories,
		"tags":        request.Tags,
		"createdAt":   request.CreatedAt,
		"updatedAt":   request.UpdatedAt,
	}

	insertResult, err := restaurantsCollection.InsertOne(ctx, filter)
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

func UpdateRestaurantAvgPrice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

	request := UpdateRestaurantAvgPriceType{}

	response := MongoJsonResponse{}

	err := c.BindJSON(&request)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(request.ID.Hex())

	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.UpdatedAt = config.GetCurrentTime()

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"avgPrice":  request.Price,
			"updatedAt": request.UpdatedAt,
		},
	}

	_, err = restaurantsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		config.Logs("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Type = "success"
	response.Data = bson.M{"ID": request.ID.Hex()}
	response.Message = "Restaurant updated"

	c.JSON(http.StatusOK, response)
}
