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

func Healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok"})
}

func PostFiles(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file[]"]
	log.Println(form)

	for _, file := range files {
		log.Println(file.Filename)
		c.SaveUploadedFile(file, "/app/resources/"+file.Filename)
	}
	c.String(http.StatusOK, "Uploaded...")
}

func GetFile(c *gin.Context) {
	name := c.Param("name")
	file := "./resources/" + name
	c.File(file)
}

func GetFileList(c *gin.Context) {
	c.String(http.StatusOK, "uploaded files list")
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

	{
		router.GET("/healthchecker", Healthcheck)
		router.POST("/upload", PostFiles)
		router.GET("/upload/:name", GetFile)
		router.GET("/upload", GetFileList)
	}

	log.Fatal(server.Run(":" + config.Port))
}
