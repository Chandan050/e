package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("cant decode the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get item form the cart")
	ErrcantBuyCartUtem    = errors.New("cannot update the perchoice")
)

func AddproductToCart() error {

}
func RemoveItemFromCart() {

}
func BuyItemFromCart() {

}
func InstabtBuyFromCart() {

}
