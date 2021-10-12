package handlers

import (
	"net/http"
	"sabino-decor/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type ProjectHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProjectHandler(ctx context.Context, collection *mongo.Collection) *ProjectHandler {
	return &ProjectHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *ProjectHandler) GetProjects(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	var projects []models.Project

	for cur.Next(handler.ctx) {
		var project models.Project
		cur.Decode(&project)

		projects = append(projects, project)
	}

	c.JSON(http.StatusOK, projects)
}

func (handler *ProjectHandler) CreateProject(c *gin.Context) {

	var project models.Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = primitive.NewObjectID()
	_, err := handler.collection.InsertOne(handler.ctx, project)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new project article"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (handler *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId},
		bson.D{
			{"$set", bson.D{
				{"name", project.Name},
				{"description", project.Description},
				{"image", project.Image},
			}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project has been updated"})
}

func (handler *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"id": objectId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project has been deleted"})
}

func (handler *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	objectid, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectid,
	})
	var article models.Project
	err := cur.Decode(&article)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}
