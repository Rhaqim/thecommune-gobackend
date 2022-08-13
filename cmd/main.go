package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/views"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/movies/", views.GetMovies).Methods("GET")
	router.HandleFunc("/movies/", views.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", views.GetMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", views.UpdateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", views.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies/", views.DeleteMovies).Methods("DELETE")
	router.HandleFunc("/restaurants/", views.Restaurants).Methods("GET")
	router.HandleFunc("/restaurant/", views.GetRestaurantByName).Methods("POST")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
