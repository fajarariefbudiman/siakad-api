package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List semua KHS beserta KRSDetail per semester
func ListKHS(c *gin.Context) {
	var khsList []models.KHS
	config.DB.Find(&khsList)

	for i := range khsList {
		var krsIDs []uint
		config.DB.Model(&models.KRS{}).Where("user_id = ? AND semester_id = ?", khsList[i].UserID, khsList[i].SemesterID).Pluck("id", &krsIDs)

		var details []models.KRSDetail
		config.DB.Preload("Course").Where("krs_id IN ?", krsIDs).Find(&details)

		khsList[i].Details = details
	}

	c.JSON(http.StatusOK, gin.H{"data": khsList})
}

// Get KHS by user_id
func GetKHSByUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userID, _ := strconv.ParseUint(userIdStr, 10, 64)

	var khsList []models.KHS
	config.DB.Where("user_id = ?", userID).Find(&khsList)

	for i := range khsList {
		var krsIDs []uint
		config.DB.Model(&models.KRS{}).Where("user_id = ? AND semester_id = ?", khsList[i].UserID, khsList[i].SemesterID).Pluck("id", &krsIDs)

		var details []models.KRSDetail
		config.DB.Preload("Course").Where("krs_id IN ?", krsIDs).Find(&details)

		khsList[i].Details = details
	}

	c.JSON(http.StatusOK, gin.H{"data": khsList})
}

// Create KHS
func CreateKHS(c *gin.Context) {
	var payload models.KHS
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ambil KRSDetail untuk semester itu
	var krsIDs []uint
	config.DB.Model(&models.KRS{}).Where("user_id = ? AND semester_id = ?", payload.UserID, payload.SemesterID).Pluck("id", &krsIDs)

	var details []models.KRSDetail
	config.DB.Preload("Course").Where("krs_id IN ?", krsIDs).Find(&details)

	payload.Details = details

	c.JSON(http.StatusCreated, gin.H{"data": payload})
}
