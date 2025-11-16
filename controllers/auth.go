package controllers

import (
	"api-siakad/config"
	"api-siakad/models"
	"api-siakad/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var payload registerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := config.DB.Where("email = ?", payload.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashed,
		Role:     "student",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role)
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

type loginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var payload loginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(payload.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role)
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func ForgotPassword(c *gin.Context) {
	type payload struct {
		Email string `json:"email" binding:"required,email"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek user
	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link sent"}) // jangan bocorkan email
		return
	}

	// Generate token random
	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role)

	// Simpan token
	config.DB.Model(&user).Update("reset_token", token)

	// TODO: KIRIM EMAIL
	c.JSON(http.StatusOK, gin.H{
		"message": "Reset token generated",
		"token":   token, // hanya untuk development!
	})
}

func ResetPassword(c *gin.Context) {
	type payload struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user dari context yang sudah di-set middleware AuthRequired
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPasswordHash(body.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is incorrect"})
		return
	}

	hashed, _ := utils.HashPassword(body.NewPassword)
	if err := config.DB.Model(&user).Update("password", hashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
