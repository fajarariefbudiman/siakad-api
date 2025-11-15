package controllers

import (
	"net/http"
	"strconv"

	"api-siakad/config"
	"api-siakad/models"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var user models.User
	if err := config.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	var user models.User
	if err := config.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
