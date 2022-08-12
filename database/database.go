package database

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoJsonResponse struct {
	Type    string   `json:"type"`
	Data    []bson.M `json:"data"`
	Message string   `json:"message"`
}

func connectMongoDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}

	uri = "mongodb://docker:mongopw@localhost:49154"

	log.Println("Connecting to MongoDB: ", uri)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func restaurants(w http.ResponseWriter, r *http.Request) {
	client := connectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database("lagos_restaurants").Collection("the_commune_test")

	var restaurants bson.M

	var restaurant []bson.M

	var response = MongoJsonResponse{}

	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	log.Fatal(err)

	log.Printf("Found multiple documents in restaurants: %v\n", cursor.Next(context.TODO()))

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&restaurants)
		log.Fatal(err)

		restaurant = append(restaurant, restaurants)

		response.Type = "success"
		response.Data = restaurant
		response.Message = "Restaurants found"
	}

	response.Type = "success"
	response.Data = restaurant
	response.Message = "Restaurants found"

	json.NewEncoder(w).Encode(response)

	defer client.Disconnect(context.TODO())

	// if err = cursor.All(context.TODO(), &restaurants); err != nil {
	// 	response.Type = "error"
	// 	response.Message = "Movie not found"
	// } else {
	// 	for cursor.Next(context.TODO()) {
	// 		var movie bson.M
	// 		err := cursor.Decode(&movie)
	// 		checkErr(err)

	// 		restaurants = append(restaurants, movie)
	// 	}

	// if r.FormValue("movie_name") == "" {
	// 	response.Type = "error"
	// 	response.Message = "Movie Name is required"
	// } else {
	// 	var movies []bson.M

	// 	filter := bson.M{"movie_name": r.FormValue("movie_name")}

	// 	cursor, err := collection.Find(context.TODO(), filter)
	// 	checkErr(err)

	// 	if err = cursor.All(context.TODO(), &movies); err != nil {
	// 		response.Type = "error"
	// 		response.Message = "Movie not found"
	// 	} else {
	// 		for cursor.Next(context.TODO()) {
	// 			var movie bson.M
	// 			err := cursor.Decode(&movie)
	// 			checkErr(err)

	// 			movies = append(movies, movie)

	// 		}

	// response.Type = "success"
	// response.Data = bson.M{"restaurants": restaurants}
	// response.Message = "Restaurants found"

	// json.NewEncoder(w).Encode(response)
}