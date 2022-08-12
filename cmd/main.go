package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_USER = "postgres"
	DB_PASS = "postgres"
	DB_NAME = "postgres"
)

type Movie struct {
	MovieID   int    `json:"movie_id"`
	MovieName string `json:"movie_name"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

type MongoJsonResponse struct {
	Type    string   `json:"type"`
	Data    []bson.M `json:"data"`
	Message string   `json:"message"`
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASS, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/movies/", getMovies).Methods("GET")
	router.HandleFunc("/movies/", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	router.HandleFunc("/movies/", deleteMovies).Methods("DELETE")
	router.HandleFunc("/restaurants", restaurants).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
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
	checkErr(err)

	log.Printf("Found multiple documents in restaurants: %v\n", cursor.Next(context.TODO()))

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&restaurants)
		checkErr(err)

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

func getMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM movies")
	checkErr(err)

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.MovieID, &movie.MovieName)
		checkErr(err)

		printMessage(fmt.Sprintf("Movie ID: %d", movie.MovieID))

		movies = append(movies, movie)
	}

	var response = JsonResponse{
		Type:    "success",
		Data:    movies,
		Message: "",
	}

	json.NewEncoder(w).Encode(response)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("movie_id")
	movieName := r.FormValue("movie_name")

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response.Type = "error"
		response.Message = "Movie ID and Movie Name are required"
	} else if movieID == "0" {
		response.Type = "error"
		response.Message = "Movie ID cannot be 0"
	} else {
		db := setupDB()
		defer db.Close()

		var insertID int

		err := db.QueryRow("INSERT INTO movies (movie_id, movie_name) VALUES ($1, $2) RETURNING movie_id", movieID, movieName).Scan(&insertID)
		checkErr(err)

		response.Type = "success"
		response.Message = "Movie created"
	}

	json.NewEncoder(w).Encode(response)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]

	var response = JsonResponse{}

	if movieID == "" {
		response.Type = "error"
		response.Message = "Movie ID is required"
	} else {
		db := setupDB()
		defer db.Close()

		var movie Movie

		err := db.QueryRow("SELECT * FROM movies WHERE movie_id = $1", movieID).Scan(&movie.MovieID, &movie.MovieName)
		checkErr(err)

		response.Type = "success"
		response.Data = append(response.Data, movie)
		response.Message = "Movie found"
	}

	json.NewEncoder(w).Encode(response)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]
	movieName := r.FormValue("movie_name")

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response.Type = "error"
		response.Message = "Movie ID and Movie Name are required"
	} else if movieID == "0" {
		response.Type = "error"
		response.Message = "Movie ID cannot be 0"
	} else {
		db := setupDB()
		defer db.Close()

		_, err := db.Exec("UPDATE movies SET movie_name = $1 WHERE movie_id = $2", movieName, movieID)
		checkErr(err)

		response.Type = "success"
		response.Message = "Movie updated"
	}

	json.NewEncoder(w).Encode(response)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]

	var response = JsonResponse{}

	if movieID == "" {
		response.Type = "error"
		response.Message = "Movie ID is required"
	} else {
		db := setupDB()
		defer db.Close()

		_, err := db.Exec("DELETE FROM movies WHERE movie_id = $1", movieID)
		checkErr(err)

		response.Type = "success"
		response.Message = "Movie deleted"
	}

	json.NewEncoder(w).Encode(response)
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM movies")
	checkErr(err)

	var response = JsonResponse{}

	response.Type = "success"
	response.Message = "Movies deleted"

	json.NewEncoder(w).Encode(response)
}
