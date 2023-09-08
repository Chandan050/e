package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Chandan050/ecommerce/database"
	"github.com/Chandan050/ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollcetion *mongo.Collection
}

func NewApplication(prodCollection, userColletion *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollcetion: userColletion,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ptoduct id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Fatalln("user ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddproductToCart(ctx, app.prodCollection, app.userCollcetion, productID, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Succesfully added to the product")

	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ptoduct id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Fatalln("user ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveItemFromCart(ctx, app.prodCollection, app.userCollcetion, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "removed from cart")

	}

}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H("error":"inavlid ID"))
			c.Abort()
			return
		}

		userAID, _ := primitive.ObjectIDFromHex(user_id)
		
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var filtedcar models.User

		err := UserCollection.FindOne(ctc, bson.D{primitive.E{Key:"_id", Value: userAID}})

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		filter_match := bson.D{{Key:"$match", Value:bson.D{primitive.E{Key:"_id",Value:userAID}}}}

		unwind := bson.D{{Key:"$unwind",Value: bson.D{primitive.E{Key:"path",Value:"$usercart"}}}}

		grouping := bson.D{{Key:"$group", Value: bson.D{primitive.E{Key:"_id",Value: "$_id"},{Key:"total", Value:bson.D{primitive.E{Key:"$sum",Value: userAID}}}}}}
		
		pointCurser, err := UserCollction.Aggregate(ctx,mongo.Pipeline{filter_match,unwind,grouping})
		
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointCurser.All(ctx,&listing); err != nil {
			log.Println(err)
			c.Abort()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing{
			c.IndentedJSON(200,json["total"])
			c.IndentedJSON(200,filtedcar.UserCart)
		}

	}

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Fatalln("user id is empty")
			c.AbortWithError(http.StatusInternalServerError, errors.New("user id is empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		err = database.BuyItemFromCart(ctx, app.userCollcetion, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully order the products")

	}
}

func (app *Application) InstabtBuy() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ptoduct id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Fatalln("user ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstabtBuyFromCart(ctx, app.prodCollection, app.userCollcetion, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "succesfully placed the ordered")
	}

}
