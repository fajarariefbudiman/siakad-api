package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateKRS(c *gin.Context) {
	var payload models.KRS
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}

func GetKRSByUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userID, _ := strconv.ParseUint(userIdStr, 10, 64)
	var krss []models.KRS
	config.DB.Preload("Details").Where("user_id = ?", userID).Find(&krss)
	c.JSON(http.StatusOK, gin.H{"data": krss})
}
