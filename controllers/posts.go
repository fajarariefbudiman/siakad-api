package controllers

import (
	"net/http"
	"strconv"

	"api-siakad/config"
	"api-siakad/models"

	"github.com/gin-gonic/gin"
)

func ListPosts(c *gin.Context) {
	var rows []models.Post
	config.DB.Order("created_at desc").Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func CreatePost(c *gin.Context) {
	var payload models.Post
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&payload)
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}

func GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Post
	if err := config.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}
