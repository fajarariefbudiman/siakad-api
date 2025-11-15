package controllers

import (
	"net/http"
	"strconv"

	"api-siakad/config"
	"api-siakad/models"

	"github.com/gin-gonic/gin"
)

func ListSemesters(c *gin.Context) {
	var items []models.Semester
	config.DB.Order("created_at desc").Find(&items)
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func CreateSemester(c *gin.Context) {
	var payload models.Semester
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&payload)
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}

func GetSemester(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Semester
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func UpdateSemester(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Semester
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&item)
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func DeleteSemester(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	config.DB.Delete(&models.Semester{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
