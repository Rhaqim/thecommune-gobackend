package views

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	restaurantCollection = "the_commune_test"
	restaurantDB         = "lagos_restaurants"
)

type MongoJsonResponse struct {
	Type    string   `json:"type"`
	Data    []bson.M `json:"data"`
	Message string   `json:"message"`
}

type Restaurant struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	Rating       string `json:"rating"`
	Reviews      string `json:"reviews"`
	OpeningHours string `json:"opening_hours"`
	Price        string `json:"price"`
	Categories   string `json:"categories"`
	Tags         string `json:"tags"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RestaurantResponse struct {
	Type    string       `json:"type"`
	Data    []Restaurant `json:"data"`
	Message string       `json:"message"`
}

type RestaurantRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RestaurantUpdateRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RestaurantDeleteRequest struct {
	ID string `json:"id"`
}

type RestaurantCreateRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	Rating       string `json:"rating"`
	Reviews      string `json:"reviews"`
	OpeningHours string `json:"opening_hours"`
	Price        string `json:"price"`
	Categories   string `json:"categories"`
	Tags         string `json:"tags"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RestaurantUpdateResponse struct {
	Type    string     `json:"type"`
	Data    Restaurant `json:"data"`
	Message string     `json:"message"`
}

type RestaurantDeleteResponse struct {
	Type    string     `json:"type"`
	Data    Restaurant `json:"data"`
	Message string     `json:"message"`
}

func Restaurants(w http.ResponseWriter, r *http.Request) {

	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

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

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

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

	collection := client.Database(restaurantDB).Collection(restaurantCollection)

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
