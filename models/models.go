package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      string             `json:"firstname" validate:"required,min=2, max=40"`
	LastName       string             `json:"lastname" validate:"required,min=1, max=40"`
	Password       string             `json:"password" validate:"required,min=6"`
	Email          string             `json:"email"  validate:"required"`
	Phone          string             `json:"phone"  validate:"required"`
	Token          string             `json:"token" `
	RefreshToken   string             `json:"refreshtoken" `
	Created_At     time.Time          `json:"created_at" `
	Updated_At     time.Time          `json:"updated_at" `
	User_ID        string             `json:"user_id" `
	UserCart       []ProductUser      `json:"usercart" bson:"usercart" `
	AddressDetails []Address          `json:"addressdetails" bson:"address"`
	OrderStatus    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName string             `json:"product_name" `
	Price       uint64             `json:"price" `
	Rating      uint8              `json:"rating" `
	Image       string             `json:"image" `
}
type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName string             `json:"product_name" bson:"product_name" `
	Price       int                `json:"price" bson:"price" `
	Rating      uint               `json:"rating" bson:"rating" `
	Image       string             `json:"image" bson:"image" `
}
type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	House     string             `json:"house_name" bson:"huose_name"`
	Street    string             `json:"street_name" bson:"street_name"`
	City      string             `json:"city_name" bson:"city_name"`
	Pincode   string             `json:"pincode" bson:"pincode"`
}
type Order struct {
	OrderId       primitive.ObjectID `bson:"_id"`
	OrderCart     []ProductUser      `json:"order_list" bson:"order_list"`
	OrderedAt     time.Time          `json:"ordered_at" bson:"ordered_at" `
	Price         int                `json:"total_price" bson:"total_price" `
	Discount      int                `json:"dicount" bson:"discount" `
	PaymentMethod Payment            `json:"pament_struct" bson:"payment_struct" `
}
type Payment struct {
	Digital bool
	COD     bool
}
