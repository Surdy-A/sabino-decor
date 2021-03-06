package handlers

import (
	"net/http"
	"sabino-decor/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type BlogHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewBlogHandler(ctx context.Context, collection *mongo.Collection) *BlogHandler {
	return &BlogHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *BlogHandler) GetArticles(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	var blogs []models.Blog

	for cur.Next(handler.ctx) {
		var blog models.Blog
		cur.Decode(&blog)

		blogs = append(blogs, blog)
	}

	c.JSON(http.StatusOK, blogs)
}

func (handler *BlogHandler) CreateArticle(c *gin.Context) {

	var blog models.Blog

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.ID = primitive.NewObjectID()
	blog.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, blog)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new blog article"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (handler *BlogHandler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId},
		bson.D{
			{"$set", bson.D{
				{"title", blog.Tittle},
				{"Conetent", blog.Content},
				{"published_at", blog.PublishedAt},
				{"image", blog.Image},
				{"comments", blog.Comments},
				{"likes", blog.Likes},
			}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog Article has been updated"})
}

func (handler *BlogHandler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"id": objectId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog Article has been deleted"})
}

func (handler *BlogHandler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	objectid, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectid,
	})
	var article models.Blog
	err := cur.Decode(&article)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}
