package handlers

import (
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/views"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	routes := Routes{
		Route{
			"Restaurants",
			"GET",
			"/restaurants",
			views.Restaurants,
		},
		Route{
			"GetRestaurantByName",
			"POST",
			"/restaurant",
			views.GetRestaurantByName,
		},
		Route{
			"GetRestaurantReviews",
			"GET",
			"/restaurant/reviews",
			views.GetRestaurantReviews,
		},
		Route{
			"UpdateRestaurantReview",
			"POST",
			"/restaurant/review",
			views.AddNewRestaurantReview,
		},
		Route{
			"UpdateReviewLikeAndDislike",
			"POST",
			"/restaurant/review/like",
			views.UpdateReviewLikeAndDislike,
		},
	}

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
