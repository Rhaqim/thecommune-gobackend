package handlers

import (
	"github.com/Rhaqim/thecommune-gobackend/pkg/views"

	"github.com/gin-gonic/gin"
)

func GinRouter() *gin.Engine {
	router := gin.Default()

	restaurants := router.Group("/restaurants")
	{
		reviews := restaurants.Group("/reviews")
		{
			reviews.GET("", views.GetRestaurantReviews)
			reviews.POST("", views.AddNewRestaurantReview)
			reviews.POST("/update", views.UpdateReviewLikeAndDislike)

		}
		restaurants.GET("", views.GetRestaurants)
		restaurants.POST("", views.CreateRestaurant)
	}

	users := router.Group("/users")
	{
		users.GET("/getuser", views.GetUserByID)
	}

	// router.POST("/restaurant", views.GetRestaurantByName)
	// router.GET("/restaurant/reviews", views.GetRestaurantReviews)
	// router.POST("/restaurant/review", views.AddNewRestaurantReview)
	// router.POST("/restaurant/review/like", views.UpdateReviewLikeAndDislike)

	return router
}

// type Route struct {
// 	Name        string
// 	Method      string
// 	Pattern     string
// 	HandlerFunc http.HandlerFunc
// }

// type Routes []Route

// func NewRouter() *mux.Router {
// 	router := mux.NewRouter()

// 	routes := Routes{
// 		Route{
// 			"Restaurants",
// 			"GET",
// 			"/restaurants",
// 			views.Restaurants,
// 		},
// 		Route{
// 			"GetRestaurantByName",
// 			"POST",
// 			"/restaurant",
// 			views.GetRestaurantByName,
// 		},
// 		Route{
// 			"GetRestaurantReviews",
// 			"GET",
// 			"/restaurant/reviews",
// 			views.GetRestaurantReviews,
// 		},
// 		Route{
// 			"UpdateRestaurantReview",
// 			"POST",
// 			"/restaurant/review",
// 			views.AddNewRestaurantReview,
// 		},
// 		Route{
// 			"UpdateReviewLikeAndDislike",
// 			"POST",
// 			"/restaurant/review/like",
// 			views.UpdateReviewLikeAndDislike,
// 		},
// 	}

// 	for _, route := range routes {
// 		router.
// 			Methods(route.Method).
// 			Path(route.Pattern).
// 			Name(route.Name).
// 			Handler(route.HandlerFunc)
// 	}

// 	return router
// }
