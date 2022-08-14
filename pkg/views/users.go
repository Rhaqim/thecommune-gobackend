package views

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUser struct {
	ID string `json:"id"`
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantCollection).Collection("USERS")

	var users bson.M
	var user []bson.M
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

		filter := bson.M{"_id": bson.M{"$ref": "USERS", "$id": id}}

		cursor, err := collection.Find(context.TODO(), filter)

		config.CheckErr(err)
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&users)
			config.CheckErr(err)

			user = append(user, users)

		}

		response.Type = "success"
		response.Data = user
		response.Message = "Users found"

		json.NewEncoder(w).Encode(response)
	}
}

type CreatUser struct {
	First_Name string             `json:"first_name"`
	Last_Name  string             `json:"last_name"`
	Username   string             `json:"username"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	Password   string             `json:"password"`
	Created_At primitive.DateTime `json:"created_at"`
	Updated_At primitive.DateTime `json:"updated_at"`
}

func CreatNewUser(w http.ResponseWriter, r *http.Request) {
	client := database.ConnectMongoDB()

	defer client.Disconnect(context.TODO())

	collection := client.Database(restaurantCollection).Collection("USERS")
	var user = CreatUser{}
	var response = MongoJsonResponse{}
	json.NewDecoder(r.Body).Decode(&user)
	log.Println("user: ", user)
	if user.Username == "" || user.Email == "" || user.Password == "" {
		response.Type = "error"
		response.Message = "username, email and password are required"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		user.Created_At = primitive.NewDateTimeFromTime(time.Now())
		user.Updated_At = primitive.NewDateTimeFromTime(time.Now())
		filter := bson.M{
			"first_name": user.First_Name,
			"last_name":  user.Last_Name,
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"created_at": user.Created_At,
			"updated_at": user.Updated_At,
		}
		insertResult, err := collection.InsertOne(context.TODO(), filter)
		config.CheckErr(err)
		log.Println("insertResult: ", insertResult)
		response.Type = "success"
		response.Message = "User created"
		json.NewEncoder(w).Encode(response)
	}
}
