package views

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	restaurantCollection = "the_commune_test"
	restaurantDB         = "lagos_restaurants"
)

type MongoJsonResponse struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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

/*
Update Review Likes and Dislikes
*/
type UpdateLikeAndDislike struct {
	ID         primitive.ObjectID `json:"id"`
	Like       int                `json:"like"`
	Dislike    int                `json:"dislike"`
	Updated_At primitive.DateTime `json:"updated_at"`
}

type Restaurant struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	Rating       string `json:"rating"`
	Reviews      string `json:"reviews"`
	OpeningHours string `json:"opening_hours"`
	Price        string `json:"price"`
	Categories   string `json:"categories"`
	Tags         string `json:"tags"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RestaurantDeleteRequest struct {
	ID string `json:"id"`
}

type RestaurantCreateRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	Rating       string `json:"rating"`
	Reviews      string `json:"reviews"`
	OpeningHours string `json:"opening_hours"`
	Price        string `json:"price"`
	Categories   string `json:"categories"`
	Tags         string `json:"tags"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RestaurantUpdateResponse struct {
	Type    string     `json:"type"`
	Data    Restaurant `json:"data"`
	Message string     `json:"message"`
}

type RestaurantDeleteResponse struct {
	Type    string     `json:"type"`
	Data    Restaurant `json:"data"`
	Message string     `json:"message"`
}

type GetRestaurantReviewsType struct {
	ID primitive.ObjectID `json:"restaurant_id"`
}