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
Get All Reviews for a Restaurant
*/

var reviewsCollection = database.OpenCollection(database.ConnectMongoDB(), DB, REVIEWS)

func GetRestaurantReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var restaurants []bson.M

	var response MongoJsonResponse

	var request = GetRestaurantReviewsType{}

	// idFromQuery := r.FormValue("id")
	if err := c.BindJSON(&request); err != nil {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}

	id, err := primitive.ObjectIDFromHex(request.ID.Hex())
	config.CheckErr(err)

	filter := bson.M{"restaurant_id": bson.M{"$ref": "RESTAURANTS", "$id": id}}

	cursor, err := reviewsCollection.Find(ctx, filter)
	if err != nil {
		response.Type = "error"
		response.Message = "Error getting restaurant reviews"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}

	if err := cursor.All(context.Background(), &restaurants); err != nil {
		config.Logs("error", err.Error())
		response.Type = "error"
		response.Message = "Error getting restaurant reviews"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}

	response.Type = "success"
	response.Data = restaurants
	response.Message = "Restaurants found"

	c.JSON(http.StatusOK, response)
}

/*
Add A New Restaurant Review
*/
func AddNewRestaurantReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

	var response = MongoJsonResponse{}
	var request = AddRestaurantReviewRequest{}
	if err := c.BindJSON(&request); err != nil {
		response.Type = "error"
		response.Message = "id is required"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}
	log.Println("request: ", request)
	restaurant_id, err := primitive.ObjectIDFromHex(request.Restaurant_ID.Hex())
	config.CheckErr(err)

	user_id, err := primitive.ObjectIDFromHex(request.Reviewer.Hex())
	config.CheckErr(err)

	request.CreatedAt = config.GetCurrentTime()
	request.UpdatedAt = config.GetCurrentTime()
	filter := bson.M{
		"reviewer":      bson.M{"$ref": "USERS", "$id": user_id},
		"review":        request.Review,
		"reviewRating":  request.Rating,
		"spent":         request.Spent,
		"reviewImages":  request.ReviewImages,
		"restaurant_id": bson.M{"$ref": "RESTAURANTS", "$id": restaurant_id},
		"dislike":       request.Dislike,
		"like":          request.Like,
		"createdAt":     request.CreatedAt,
		"updatedAt":     request.UpdatedAt,
	}
	insertResult, err := reviewsCollection.InsertOne(ctx, filter)
	if err != nil {
		response.Type = "error"
		response.Message = "Error adding new restaurant review"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}
	log.Println("insertResult: ", insertResult)
	response.Type = "success"
	response.Data = bson.M{"InsertID": insertResult.InsertedID.(primitive.ObjectID).Hex()}
	response.Message = "Review added"
	c.JSON(http.StatusOK, response)
}

/*
Update Review Likes and Dislikes
*/
func UpdateReviewLikeAndDislike(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.ConnectMongoDB().Disconnect(context.TODO())

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
	id, err := primitive.ObjectIDFromHex(request.ID.Hex())
	config.CheckErr(err)
	request.UpdatedAt = config.GetCurrentTime()
	updateResult, err := reviewsCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"like": request.Like, "dislike": request.Dislike, "updatedAt": request.UpdatedAt}})
	if err != nil {
		response.Type = "error"
		response.Message = "Error updating review like and dislike"
		c.JSON(http.StatusBadRequest, gin.H{"res": response})
		return
	}
	log.Println("updateResult: ", updateResult)
	response.Type = "success"
	response.Data = bson.M{"UpdateID": updateResult}
	response.Message = "Review updated"
	c.JSON(http.StatusOK, response)
}
