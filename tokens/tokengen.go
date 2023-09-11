package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Chandan050/ecommerce/database"
	jwt "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SigneDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedtoken string, signedrefreshtoken string, err error) {
	claims := &SigneDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SigneDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return token, refreshtoken, err

}

func ValidateToken(sinedtoken string)(claims *SigneDetails, msg string) {
	token , err := jwt.ParseWithClaims(sinedtoken,&SigneDetails{}, func(t *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "Signing Method ES256"
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
	}
	claims, ok := token.Claims.(*SigneDetails)
	if !ok ||!token.Valid {
		msg= "the token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}
	return claims,msg
} 

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {

	var ctx, cancel := context.WithTimeout(context.Background(), 100* time.Second)

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})

	upsert := true

	filter := bson.M("user_id":userid)

	opt := options.UpdateOptions{
		Upsert : &upsert,
	}

	UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj}
	},
	&opt)

defer cancel()
if err != nil {
	log.Panic(err)
	return
}

}
