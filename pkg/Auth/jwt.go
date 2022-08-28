package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Rhaqim/thecommune-gobackend/pkg/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var jwtKey = []byte("supersecretkey")

var collection *mongo.Collection = database.OpenCollection(database.ConnectMongoDB(), "lagos_restaurants", "USERS")

type JWTClaim struct {
	Username string             `json:"username"`
	Email    string             `json:"email"`
	UserId   primitive.ObjectID `json:"userId"`
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET")

func GenerateJWT(email string, username string, userid primitive.ObjectID) (token string, refreshToken string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		UserId:   userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshClaims := &JWTClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return
}
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte([]byte(SECRET_KEY)), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func UpdateToken(signedToken string, signedRefreshToken string, email string, username string, userid primitive.ObjectID) (token string, refreshToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updated_at})

	upsert := true

	userid, err = primitive.ObjectIDFromHex(userid.Hex())

	filter := bson.M{"_id": userid}
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opts)
	if err != nil {
		log.Panic(err)
		return
	}

	token, refreshToken, err = GenerateJWT(email, username, userid)
	if err != nil {
		return
	}

	return
}
