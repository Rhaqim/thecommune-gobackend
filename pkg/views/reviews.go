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

/*
Add A New Restaurant Review
*/

type AddRestaurantReviewRequest struct {
	Reviewer primitive.ObjectID `json:"reviewer"`

	Review string `json:"review"`

	Rating int `json:"rating"`

	Spent int `json:"spent"`

	Review_Images []string `json:"review_images"`

	Restaurant_ID primitive.ObjectID `json:"restaurant_id"`

	Dislike int `json:"dislike"`

	Like int `json:"like"`

	Created_At primitive.DateTime `json:"created_at"`

	Updated_At primitive.DateTime `json:"updated_at"`
}

func AddNewRestaurantReview(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")
	var response = MongoJsonResponse{}
	var request = AddRestaurantReviewRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	config.CheckErr(err)
	log.Println("request: ", request)
	if request.Review == "" || request.Restaurant_ID == primitive.NilObjectID || request.Reviewer == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "name, review, rating, restaurant_id, user_id are required"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		restaurant_id, err := primitive.ObjectIDFromHex(request.Restaurant_ID.Hex())
		config.CheckErr(err)

		user_id, err := primitive.ObjectIDFromHex(request.Reviewer.Hex())
		config.CheckErr(err)

		request.Created_At = config.GetCurrentTime()
		request.Updated_At = config.GetCurrentTime()
		filter := bson.M{
			"reviewer":      bson.M{"$ref": "USERS", "$id": user_id},
			"review":        request.Review,
			"rating":        request.Rating,
			"spent":         request.Spent,
			"review_images": request.Review_Images,
			"restaurant_id": bson.M{"$ref": "RESTAURANTS", "$id": restaurant_id},
			"dislike":       request.Dislike,
			"like":          request.Like,
			"created_at":    request.Created_At,
			"updated_at":    request.Updated_At,
		}
		insertResult, err := collection.InsertOne(context.TODO(), filter)
		config.CheckErr(err)
		log.Println("insertResult: ", insertResult)
		response.Type = "success"
		response.SingleData = bson.M{"InsertID": insertResult.InsertedID.(primitive.ObjectID).Hex()}
		response.Message = "Review added"
		json.NewEncoder(w).Encode(response)
	}
}

/*
Increment Like Count
*/
type UpdateLikeAndDislike struct {
	ID         primitive.ObjectID `json:"id"`
	Like       int                `json:"like"`
	Dislike    int                `json:"dislike"`
	Updated_At primitive.DateTime `json:"updated_at"`
}

func UpdateReviewLikeAndDislike(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")
	var response = MongoJsonResponse{}
	var request = UpdateLikeAndDislike{}
	err := json.NewDecoder(r.Body).Decode(&request)
	config.CheckErr(err)
	log.Println("request: ", request)
	if request.ID == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "id is required"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		id, err := primitive.ObjectIDFromHex(request.ID.Hex())
		config.CheckErr(err)
		request.Updated_At = config.GetCurrentTime()
		updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"like": request.Like, "dislike": request.Dislike, "updated_at": request.Updated_At}})
		config.CheckErr(err)
		log.Println("updateResult: ", updateResult)
		response.Type = "success"
		response.SingleData = bson.M{"UpdateID": updateResult}
		response.Message = "Review updated"
		json.NewEncoder(w).Encode(response)
	}
}
