package main

import (
	"log"
	"os"

	"github.com/Chandan050/ecommerce/controllers"
	"github.com/Chandan050/ecommerce/database"
	"github.com/Chandan050/ecommerce/middleware"
	"github.com/Chandan050/ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantBuy", app.InstabtBuy())

	log.Fatal(router.Run(":" + port))

}
