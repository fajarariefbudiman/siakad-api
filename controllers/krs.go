package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateKRS(c *gin.Context) {
	var payload struct {
		UserID     uint            `json:"user_id"`
		SemesterID uint            `json:"semester_id"`
		Finalized  bool            `json:"finalized"`
		Courses    []models.Course `json:"courses"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah KRS user + semester sudah ada
	var existingKRS models.KRS
	if err := config.DB.Where("user_id = ? AND semester_id = ?", payload.UserID, payload.SemesterID).First(&existingKRS).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "KRS untuk semester ini sudah ada"})
		return
	}

	// Buat KRS baru
	krs := models.KRS{
		UserID:     payload.UserID,
		SemesterID: payload.SemesterID,
		Finalized:  payload.Finalized,
	}
	if err := config.DB.Create(&krs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tambahkan Courses
	for i := range payload.Courses {
		payload.Courses[i].KRSID = krs.ID
	}
	if len(payload.Courses) > 0 {
		if err := config.DB.Create(&payload.Courses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Ambil KRS beserta courses yang baru dibuat
	var createdKRS models.KRS
	config.DB.Preload("Courses").First(&createdKRS, krs.ID)

	c.JSON(http.StatusCreated, gin.H{"data": createdKRS})
}

func GetKRSByUser(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userID, _ := strconv.ParseUint(userIdStr, 10, 64)

	var krss []models.KRS
	// preload Courses, bukan Details
	config.DB.Preload("Courses").Where("user_id = ?", userID).Find(&krss)

	c.JSON(http.StatusOK, gin.H{"data": krss})
}
