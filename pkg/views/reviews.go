package views

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRestaurantReviews(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")

	var restaurants bson.M

	var restaurant []bson.M
	var response = MongoJsonResponse{}

	idFromQuery := r.FormValue("id")
	log.Println("id: ", idFromQuery)

	if idFromQuery == "" {
		response.Type = "error"
		response.Message = "id is required"
		json.NewEncoder(w).Encode(response)
		return
	} else {

		id, err := primitive.ObjectIDFromHex(idFromQuery)
		config.CheckErr(err)

		filter := bson.M{"restaurant_id": bson.M{"$ref": "RESTAURANTS", "$id": id}}

		cursor, err := collection.Find(context.TODO(), filter)
		config.CheckErr(err)

		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&restaurants)
			config.CheckErr(err)

			restaurant = append(restaurant, restaurants)

		}

		response.Type = "success"
		response.Data = restaurant
		response.Message = "Restaurants found"

		json.NewEncoder(w).Encode(response)
	}
}

type UpdateRestaurantReviewRequest struct {
	Name string `json:"name"`

	Review string `json:"review"`

	Rating string `json:"rating"`

	Restaurant_ID primitive.ObjectID `json:"restaurant_id"`

	Dislike int `json:"dislike"`

	Like int `json:"like"`

	Spent int `json:"spent"`

	ReviewImages []string `json:"review_images"`

	Reviewer primitive.ObjectID `json:"user_id"`

	Created_At primitive.DateTime `json:"created_at"`

	Updated_At primitive.DateTime `json:"updated_at"`
}

func AddNewRestaurantReview(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")
	var response = MongoJsonResponse{}
	var request = UpdateRestaurantReviewRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	config.CheckErr(err)
	log.Println("request: ", request)
	if request.Name == "" || request.Review == "" || request.Rating == "" || request.Restaurant_ID == primitive.NilObjectID || request.Reviewer == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "name, review, rating, restaurant_id, user_id are required"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		id, err := primitive.ObjectIDFromHex(request.Restaurant_ID.Hex())
		config.CheckErr(err)
		request.Restaurant_ID = id
		request.Created_At = config.GetCurrentTime()
		request.Updated_At = config.GetCurrentTime()
		insertResult, err := collection.InsertOne(context.TODO(), request)
		config.CheckErr(err)
		log.Println("insertResult: ", insertResult)
		response.Type = "success"
		response.SingleData = bson.M{insertResult.InsertedID.(primitive.ObjectID).Hex(): request}
		response.Message = "Review added"
		json.NewEncoder(w).Encode(response)
	}
}
