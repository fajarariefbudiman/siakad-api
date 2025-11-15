package main

import (
	"api-siakad/config"
	"api-siakad/controllers"
	"api-siakad/middleware"
	"api-siakad/models"
	"api-siakad/seeders"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to environment variables")
	}

	// connect DB
	config.ConnectDatabase()

	// Auto migrate
	db := config.DB
	db.AutoMigrate(&models.User{}, &models.Semester{}, &models.KRS{}, &models.KRSDetail{}, &models.KHS{}, &models.Payment{}, &models.Post{})

	seeders.Seed()

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/auth/forgot-password", controllers.ForgotPassword)

		// protected
		auth := api.Group("/")
		auth.Use(middleware.AuthRequired())
		{
			auth.GET("/me", controllers.Me)
			// users
			auth.GET("/users/:id", controllers.GetUserByID)
			auth.POST("/reset-password", controllers.ResetPassword)

			// semester
			auth.GET("/semester", controllers.ListSemesters)
			auth.POST("/semester", controllers.CreateSemester)
			auth.GET("/semester/:id", controllers.GetSemester)
			auth.PUT("/semester/:id", controllers.UpdateSemester)
			auth.DELETE("/semester/:id", controllers.DeleteSemester)

			// krs
			auth.POST("/krs", controllers.CreateKRS)
			auth.GET("/krs/user/:user_id", controllers.GetKRSByUser)

			// khs
			auth.GET("/khs", controllers.ListKHS)
			auth.POST("/khs", controllers.CreateKHS)

			// payments
			auth.GET("/payments", controllers.ListPayments)
			auth.POST("/payments", controllers.CreatePayment)
			auth.GET("/payments/user/:user_id", controllers.GetPaymentsByUser)

			// posts
			auth.GET("/posts", controllers.ListPosts)
			auth.POST("/posts", controllers.CreatePost)
			auth.GET("/posts/:id", controllers.GetPost)
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
