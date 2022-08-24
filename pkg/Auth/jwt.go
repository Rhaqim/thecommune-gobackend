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

var jwtKey = []byte("supersecretkey")

var collection *mongo.Collection = database.OpenCollection(database.ConnectMongoDB(), "lagos_restaurants", "USERS")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET")

func GenerateJWT(email string, username string) (token string, refreshToken string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshClaims := &JWTClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtKey)
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
			return []byte(jwtKey), nil
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

func UpdateToken(signedToken string, signedRefreshToken string, email string, username string) (token string, refreshToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, primitive.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, primitive.E{Key: "refreshToken", Value: signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updatedAt", Value: updated_at})

	upsert := true

	filter := bson.D{{"email", email}}
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = collection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opts)
	if err != nil {
		log.Fatal(err)
		return
	}

	token, refreshToken, err = GenerateJWT(email, username)
	if err != nil {
		return
	}

	return
}
