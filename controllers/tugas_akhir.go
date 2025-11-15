package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTugasAkhir(c *gin.Context) {
	var input models.TugasAkhir

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check user exists
	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Save
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tugas Akhir created successfully",
		"data":    input,
	})
}

func GetTugasAkhirByCategory(c *gin.Context) {
	category := c.Param("category")

	var list []models.TugasAkhir

	if err := config.DB.Where("category = ?", category).
		Preload("User").
		Find(&list).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
