package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListKHS(c *gin.Context) {
	var rows []models.KHS
	config.DB.Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func CreateKHS(c *gin.Context) {
	var payload models.KHS
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&payload)
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}
