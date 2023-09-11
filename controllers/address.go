package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Chandan050/ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
			return
		}
		addressID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}

		var addresses models.Address
		addresses.AddressID = primitive.NewObjectID()
		if err := c.BindHeader(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel.Close()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: addressID}}}}

		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addressID"}}}}

		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$addressID"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCurser, err := UserCollction.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(500,"internal server error")
		}

		var addressInfo []bson.M
		if err := pointCurser.All(ctx, %adaddressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, addres_no := range addressInfo{
			count := addres_no["count"]
			size = count.(int32)
			if size < 2 {
				filter := bson.D{primitive.E{Key: "_id", Value: addressID}}
				update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address",Value: addresses}}}}
				_,err := UserCollction.UpdateOne(ctx, filter, update)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				c.IndentedJSON(400, "Not allowed")
			}
			
		}
		defer cancel()
		ctx.Done()
	}
}

func EditAddress() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
			return
		}
		userAID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}
		var editaddress models.Address
		if err = c.BindJSON(&editaddress)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userAID}}

		update := bson.D{{Key:"$set" , Value: bson.D{primitive.E{Key: "address.O.house_name", Value: editaddress.House}, {Key: "address.0.Street_name", Value: editaddress.Street},{Key:"address.0.city_name",Value: editaddress.City},{Key: "address.0.pin_code", Value: editaddress.Pincode}}}}

		_,err := UserCollction.UpdateOne(ctx, filter,update)
		if err != nil {
			c.IndentedJSON(500, "something went wrong")
		}
		defer cancel()
		c.IndentedJSON(200,"successfully updated")

	}

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		user_id := c.Query("id")
		if user_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
			return
		}

		addressUser := make([]models.Address, 0)

		userAID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userAID}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addressUser}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		defer cancel()

		if err != nil {
			c.IndentedJSON(404, "wrong cammand")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Sucessfully Deleted")
	}

}
