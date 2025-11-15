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
	"gorm.io/gorm"
)

// DropAllTables akan menghapus semua tabel di database
func DropAllTables(db *gorm.DB) {
	var tables []string

	// Cek jenis DB (Postgres atau MySQL/SQLite)
	dialect := db.Dialector.Name()

	if dialect == "postgres" {
		db.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public'").Scan(&tables)
	} else {
		db.Raw("SHOW TABLES").Scan(&tables)
	}

	for _, table := range tables {
		log.Printf("Dropping table: %s", table)
		db.Migrator().DropTable(table)
	}

	log.Println("All tables dropped!")
}

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to environment variables")
	}

	// connect DB
	config.ConnectDatabase()
	db := config.DB

	// Hapus semua tabel
	DropAllTables(db)

	// Migrasi ulang semua tabel
	db.AutoMigrate(&models.User{}, &models.Semester{}, &models.KRS{}, &models.Course{}, &models.KHS{}, &models.Payment{}, &models.Post{}, &models.KRSDetail{}, &models.TugasAkhir{})

	// Seed data
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
			auth.DELETE("/krs/:krs_id/course/:detail_id", controllers.DeleteCourseFromKRS)
			auth.POST("/krs/:id/course", controllers.AddCourseToKRS)

			auth.POST("/tugas-akhir", controllers.CreateTugasAkhir)
			auth.GET("/tugas-akhir/:category", controllers.GetTugasAkhirByCategory)

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
