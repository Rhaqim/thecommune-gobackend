package views

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/config"
	"github.com/Rhaqim/thecommune-gobackend/database"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoJsonResponse struct {
	Type    string   `json:"type"`
	Data    []bson.M `json:"data"`
	Message string   `json:"message"`
}

func Restaurants(w http.ResponseWriter, r *http.Request) {

	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database("lagos_restaurants").Collection("the_commune_test")

	var restaurants bson.M

	var restaurant []bson.M

	var response = MongoJsonResponse{}

	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	config.CheckErr(err)

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&restaurants)
		config.CheckErr(err)

		restaurant = append(restaurant, restaurants)

		response.Type = "success"
		response.Data = restaurant
		response.Message = "Restaurants found"
	}

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	json.NewEncoder(w).Encode(response)
}

func GetRestaurantByName(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database("lagos_restaurants").Collection("the_commune_test")

	var restaurants bson.M

	var restaurant []bson.M

	var response = MongoJsonResponse{}

	filter := bson.M{"name": "Shiro Lagos"}

	err := collection.FindOne(context.TODO(), filter).Decode(&restaurants)
	config.CheckErr(err)

	restaurant = append(restaurant, restaurants)

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	json.NewEncoder(w).Encode(response)
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database("lagos_restaurants").Collection("the_commune_test")

	var restaurants bson.M

	var restaurant []bson.M

	var response = MongoJsonResponse{}

	InsertResult, err := collection.InsertOne(context.TODO(), restaurants)

	config.CheckErr(err)

	log.Printf("Inserted a single document: %v", InsertResult.InsertedID)

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	json.NewEncoder(w).Encode(response)

}
