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
