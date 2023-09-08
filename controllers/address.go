package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Chandan050/ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

}

func EditAddress() gin.HandlerFunc {

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
