package main

// ? Require the packages
import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/circlepen/webtest/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ? Create required variables that we'll re-assign later
var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
)

func init() {

	// ? Load the .env variables
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// ? Create the Gin Engine instance
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := server.Group("/api")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/gopher", func(c *gin.Context) {
		c.File("./resources/gopher.png")
	})

	router.GET("/myname", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "circlepen",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["file[]"]
		log.Println(form)

		for _, file := range files {
			log.Println(file.Filename)
			c.SaveUploadedFile(file, "/app/resources/"+file.Filename)
		}
		c.String(http.StatusOK, "Uploaded...")
	})

	router.GET("/upload/:name", func(c *gin.Context) {
		name := c.Param("name")
		file := "./resources/" + name
		c.File(file)
	})

	router.GET("/upload", func(c *gin.Context) {
		c.String(http.StatusOK, "uploaded files list")
	})

	log.Fatal(server.Run(":" + config.Port))
}
