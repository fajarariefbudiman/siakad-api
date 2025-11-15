package controllers

import (
	"net/http"
	"strconv"

	"api-siakad/config"
	"api-siakad/models"

	"github.com/gin-gonic/gin"
)

func ListPayments(c *gin.Context) {
	var rows []models.Payment
	config.DB.Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func CreatePayment(c *gin.Context) {
	var payload models.Payment
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&payload)
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}

func GetPaymentsByUser(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	var rows []models.Payment
	config.DB.Where("user_id = ?", uid).Find(&rows)
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
