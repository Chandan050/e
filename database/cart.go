package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Chandan050/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("cant decode the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get item form the cart")
	ErrcantBuyCartUtem    = errors.New("cannot update the perchoice")
)

func AddproductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {
	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	err = searchFromDb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

}
func RemoveItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productId}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil

}
func BuyItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, userID string) error {

	//fetch the cart of the user
	//find the cart total
	// create an order the items
	//added order to the user collection
	//added items in the cart to order list
	// empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderId = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}

	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currrentResult, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if err := currrentResult.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32

	for _, item := range getusercart {
		price := item["total"]
		total_price = price.(int32)
	}

	orderCart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter1 := bson.D{primitive.E{Key: "_id", Value: id}}
	update1 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateMany(ctx, filter1, update1)
	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser, 0)
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}

	userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrcantBuyCartUtem
	}
	return nil

}

func InstabtBuyFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var order_details models.Order

	order_details.OrderId = primitive.NewObjectID()
	order_details.OrderedAt = time.Now()
	order_details.OrderCart = make([]models.ProductUser, 0)
	order_details.PaymentMethod.COD = true

	if err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details); err != nil {
		log.Panicln(err)
	}

	order_details.Price = product_details.Price

	filter := bson.D{primitive.E{Key:"_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_details}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Panicln(err)
	}

	filter2:= 

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Panicln(err)
	}
	return nil
	

}
