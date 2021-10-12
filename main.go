package main

import (
	"context"
	"log"
	"sabino-decor/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var authHandler *handlers.AuthHandler
var blogHandler *handlers.BlogHandler
var projectHandler *handlers.ProjectHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	collectionUsers := client.Database("recipedb").Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

	blogCollection := client.Database("sabinodecor").Collection("blog")
	blogHandler = handlers.NewBlogHandler(ctx, blogCollection)

	projectCollection := client.Database("sabinodecor").Collection("project")
	projectHandler = handlers.NewProjectHandler(ctx, projectCollection)
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/create", authHandler.CreateHandler)

	router.GET("articles/", blogHandler.GetArticles)
	router.GET("articles/:id", blogHandler.GetArticle)

	//Project Endpoints
	router.GET("projects/", projectHandler.GetProjects)
	router.GET("projects/:id", projectHandler.GetProject)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		router.POST("/articles", blogHandler.CreateArticle)
		router.PUT("/articles/:id", blogHandler.UpdateArticle)
		router.DELETE("/articles/:id", blogHandler.DeleteArticle)

		//Project Endpoints
		router.POST("/projects", projectHandler.CreateProject)
		router.PUT("/projects/:id", projectHandler.UpdateProject)
		router.DELETE("/projects/:id", projectHandler.DeleteProject)
	}

	router.Run()
}
