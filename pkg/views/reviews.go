package views

import (
	"context"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Get All Reviews for a Restaurant
*/
func GetRestaurantReviews(c *gin.Context) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")

	var restaurants bson.M
	var restaurant []bson.M

	var response MongoJsonResponse

	var request = GetRestaurantReviewsType{}

	// idFromQuery := r.FormValue("id")
	if err := c.BindJSON(&request); err != nil {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}

	if request.ID == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, response)
		return
	} else {

		id, err := primitive.ObjectIDFromHex(request.ID.Hex())
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

		c.JSON(http.StatusOK, response)
	}
}

/*
Add A New Restaurant Review
*/
func AddNewRestaurantReview(c *gin.Context) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")
	var response = MongoJsonResponse{}
	var request = AddRestaurantReviewRequest{}
	if err := c.BindJSON(&request); err != nil {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}
	log.Println("request: ", request)

	if request.Review == "" || request.Restaurant_ID == primitive.NilObjectID || request.Reviewer == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "name, review, rating, restaurant_id, user_id are required"
		c.JSON(http.StatusBadRequest, response)
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
		response.Data = bson.M{"InsertID": insertResult.InsertedID.(primitive.ObjectID).Hex()}
		response.Message = "Review added"
		c.JSON(http.StatusOK, response)
	}
}

/*
Update Review Likes and Dislikes
*/
func UpdateReviewLikeAndDislike(c *gin.Context) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantDB).Collection("REVIEWS")
	var response = MongoJsonResponse{}
	var request = UpdateLikeAndDislike{}
	if err := c.BindJSON(&request); err != nil {
		config.CheckErr(err)
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}
	log.Println("request: ", request)
	if request.ID == primitive.NilObjectID {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		id, err := primitive.ObjectIDFromHex(request.ID.Hex())
		config.CheckErr(err)
		request.Updated_At = config.GetCurrentTime()
		updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"like": request.Like, "dislike": request.Dislike, "updated_at": request.Updated_At}})
		config.CheckErr(err)
		log.Println("updateResult: ", updateResult)
		response.Type = "success"
		response.Data = bson.M{"UpdateID": updateResult}
		response.Message = "Review updated"
		c.JSON(http.StatusOK, response)
	}
}
