package views

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/config"
	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/gorilla/mux"
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

func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := database.SetupDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM movies")
	config.CheckErr(err)

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.MovieID, &movie.MovieName)
		config.CheckErr(err)

		config.PrintMessage(fmt.Sprintf("Movie ID: %d", movie.MovieID))

		movies = append(movies, movie)
	}

	var response = JsonResponse{
		Type:    "success",
		Data:    movies,
		Message: "",
	}

	json.NewEncoder(w).Encode(response)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
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
		db := database.SetupDB()
		defer db.Close()

		var insertID int

		err := db.QueryRow("INSERT INTO movies (movie_id, movie_name) VALUES ($1, $2) RETURNING movie_id", movieID, movieName).Scan(&insertID)
		config.CheckErr(err)

		response.Type = "success"
		response.Message = "Movie created"
	}

	json.NewEncoder(w).Encode(response)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]

	var response = JsonResponse{}

	if movieID == "" {
		response.Type = "error"
		response.Message = "Movie ID is required"
	} else {
		db := database.SetupDB()
		defer db.Close()

		var movie Movie

		err := db.QueryRow("SELECT * FROM movies WHERE movie_id = $1", movieID).Scan(&movie.MovieID, &movie.MovieName)
		config.CheckErr(err)

		response.Type = "success"
		response.Data = append(response.Data, movie)
		response.Message = "Movie found"
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
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
		db := database.SetupDB()
		defer db.Close()

		_, err := db.Exec("UPDATE movies SET movie_name = $1 WHERE movie_id = $2", movieName, movieID)
		config.CheckErr(err)

		response.Type = "success"
		response.Message = "Movie updated"
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["movie_id"]

	var response = JsonResponse{}

	if movieID == "" {
		response.Type = "error"
		response.Message = "Movie ID is required"
	} else {
		db := database.SetupDB()
		defer db.Close()

		_, err := db.Exec("DELETE FROM movies WHERE movie_id = $1", movieID)
		config.CheckErr(err)

		response.Type = "success"
		response.Message = "Movie deleted"
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := database.SetupDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM movies")
	config.CheckErr(err)

	var response = JsonResponse{}

	response.Type = "success"
	response.Message = "Movies deleted"

	json.NewEncoder(w).Encode(response)
}
