package controllers

import (
	"bytes"
	"context"
	"go/token"
	"log"
	"net/http"
	"time"

	"github.com/Chandan050/ecommerce/database"
	"github.com/Chandan050/ecommerce/models"
	"github.com/Chandan050/ecommerce/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/validate"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
var UserCollction *mongo.Collection = database.UserData(database.Client, "Users")
var productCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var validate = validator.New()

func HashPassword(password string)string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)

}

func VarifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg=  "Login or Password is incorrect"
	}
	return valid, msg


}


func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error": validationErr})
		}
		var UserCollection *mongo.Collection

		count, err := UserCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H("error":err))
			return
		}
		if count >0 {
			c.JSON(http.StatusBadRequest,gin.H("error":"user already exits"))
		}
		
		count, err := UserCollection.CountDocuments(ctx,bson.M{"email":user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H("error":err))
			return
		}
		if count >0 {
			c.JSON(http.StatusBadRequest,gin.H("error":"this user phone number already exits"))
			return
		}

		password := HashPassword(user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID= primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken,_ := generate.TokenGenerator(user.Email, user.FirstName, user.LastName,user.User_ID)

		user.Token = token
		user.RefreshToken = refreshToken
		user.UserCart = make([]models.ProductUser,0 )
		User.AddressDetails = make([]models.Address,0)
		user.OrderStatus = make([]models.Order,0)
		_, inserter := UserCollection.InsertOne(ctx, user)
		if inserter != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
			return
		}

		defer cancel()
		c.json(http.StatusCreated,"Successfully signed in!")

	}
}
func  Login() gin.HandlerFunc {
	return func(c *gin.context){
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err})
			return
		}
		var UserCollction *mongo.Collection
		var foundUser models.User
		UserCollction.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":"login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VarifyPassword(user.Password, foundUser.Password)

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		token , refreshtoken,_ :=  tokens.TokenGenerator(foundUser.Email,foundUser.FirstName,foundUser.LastName,foundUser.User_ID)

		defer cancel()

		tokens.UpdateAllTokens(token,refreshtoken,foundUser.User_ID)

		c.JSON(http.StatusFound,foundUser)	
	}
 	
}



func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func (c *gin.Context)  {

		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100* time.Second)

		defer cancel()
		cursor, err := productCollection.Find(ctx,bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong plz try again")
			return
		}
		err := cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError, "something went wrong, please try again")
			return 
		}
		defer cursor.Close()
		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return 
		}
		defer cancel()
		c.IndentedJSON(200, productList)

	}

}

func SearchProductByQuery() gin.HHandlerFunc {
	return func (c *gin.Context)  {
		c.Header("Content-Type", "application/json")
		var searchProduct []models.Product
		queyParam := c.Query("name")

		//you want to check if its empty
		if queyParam == "" {
			log.Println("query not found")
			
			c.JSON(http.StatusNotFound, gin.H("Errror":"Invalid search index"))
			c.Abort()
			return
		}
		var ctc, cancel = context.WithTimeout(contect.Background(),100*time.Second)

		defer cancel()

		cursorDb, err := productCollection.Find(ctx, bson.M{"product_name":bson.M{"$regex":queyParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong while fetching the data")
		}

		err := cursorDb.All(ctx, &searchProduct)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400,"inavlid")
			return
		}
		defer cursorDb.Close(ctx)

		if err := cursorDb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid requsest")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchProduct)
	}

}
