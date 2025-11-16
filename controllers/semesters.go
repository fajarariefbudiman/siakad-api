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

func GetSemester(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Semester
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
