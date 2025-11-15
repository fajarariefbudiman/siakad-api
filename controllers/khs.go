package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List semua KHS beserta Courses per semester
func ListKHS(c *gin.Context) {
	var khsList []models.KHS
	config.DB.Find(&khsList)

	for i := range khsList {
		var krs models.KRS
		// Ambil KRS user + semester dan preload Courses
		if err := config.DB.Preload("Courses").Where("user_id = ? AND semester_id = ?", khsList[i].UserID, khsList[i].SemesterID).First(&krs).Error; err == nil {
			khsList[i].Courses = krs.Courses
		} else {
			khsList[i].Courses = []models.Course{}
		}
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
		var krs models.KRS
		if err := config.DB.Preload("Courses").Where("user_id = ? AND semester_id = ?", khsList[i].UserID, khsList[i].SemesterID).First(&krs).Error; err == nil {
			khsList[i].Courses = krs.Courses
		} else {
			khsList[i].Courses = []models.Course{}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": khsList})
}

// Create KHS (GPA)
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

	// Ambil courses untuk semester itu
	var krs models.KRS
	if err := config.DB.Preload("Courses").Where("user_id = ? AND semester_id = ?", payload.UserID, payload.SemesterID).First(&krs).Error; err == nil {
		payload.Courses = krs.Courses
	} else {
		payload.Courses = []models.Course{}
	}

	c.JSON(http.StatusCreated, gin.H{"data": payload})
}
