package views

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	RESTAURANTS = "RESTAURANTS"
	REVIEWS     = "REVIEWS"
	USERS       = "USERS"
	DB          = "lagos_restaurants"
)

type MongoJsonResponse struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

///////////////////////////////////////////////////////////////////////////////////////
// REVIEWS
/*
Get Restaurants Reviews
*/
type GetRestaurantReviewsType struct {
	ID primitive.ObjectID `json:"restaurant_id"`
}

/*
Add A New Restaurant Review
*/

type AddRestaurantReviewRequest struct {
	Reviewer primitive.ObjectID `json:"reviewer"`

	Review string `json:"review"`

	Rating int `json:"reviewRating"`

	Spent float64 `json:"spent"`

	ReviewImages []interface{} `json:"reviewImages"`

	Restaurant_ID primitive.ObjectID `json:"restaurant_id"`

	Dislike int `json:"dislike"`

	Like int `json:"like"`

	CreatedAt primitive.DateTime `json:"createdAt"`

	UpdatedAt primitive.DateTime `json:"updatedAt"`
}

/*
Update Review Likes and Dislikes
*/
type UpdateLikeAndDislike struct {
	ID        primitive.ObjectID `json:"id"`
	Like      int                `json:"like"`
	Dislike   int                `json:"dislike"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
}

///////////////////////////////////////////////////////////////////////////////////////
// USERS
/*
Get User by ID
*/
type GetUser struct {
	ID primitive.ObjectID `json:"user_id"`
}

type CreatUser struct {
	Fullname   string             `json:"fullname"`
	Username   string             `json:"username"`
	Avatar     []interface{}      `json:"avatar"`
	Email      string             `json:"email"`
	Password   string             `json:"password"`
	Social     []interface{}      `json:"social"`
	Role       string             `json:"role"`
	Created_At primitive.DateTime `json:"created_at"`
	Updated_At primitive.DateTime `json:"updated_at"`
}

type CreateRestaurants struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Slug        string             `json:"slug"`
	Images      []interface{}      `json:"images"`
	Address     string             `json:"address"`
	Phone       string             `json:"phone"`
	Email       string             `json:"email"`
	Website     string             `json:"website"`
	Lat         string             `json:"lat"`
	Long        string             `json:"lng"`
	Rating      int                `json:"rating"`
	OpeningTime []interface{}      `json:"openingTime"`
	Currency    string             `json:"currency"`
	Price       float64            `json:"price"`
	Categories  []string           `json:"categories"`
	Tags        []string           `json:"tags"`
	CreatedAt   primitive.DateTime `json:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt"`
}

type UpdateRestaurantAvgPriceType struct {
	ID        primitive.ObjectID `json:"id"`
	Price     float64            `json:"price"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
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
