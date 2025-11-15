package seeders

import (
	"api-siakad/config"
	"api-siakad/models"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	db := config.DB
	log.Println("Running Seeder...")

	rand.Seed(time.Now().UnixNano())

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := []models.User{
		{Name: "Rahman Fajar Banyu Adji", Email: "rahman@student.com", Nim: "230441001", Password: string(hash), Role: "student"},
		{Name: "Dr. Budi Santoso", Email: "budi@dosen.com", Nim: "DSN001", Password: string(hash), Role: "dosen"},
		{Name: "Mila Rahma, M.Kom", Email: "mila@dosen.com", Nim: "DSN002", Password: string(hash), Role: "dosen"},
	}

	for _, u := range users {
		var exists models.User
		if err := db.Where("email = ?", u.Email).First(&exists).Error; err == nil {
			log.Printf("User '%s' exists, skipping...", u.Email)
		} else {
			db.Create(&u)
			log.Printf("User '%s' created", u.Email)
		}
	}

	var student models.User
	db.Where("email = ?", "rahman@student.com").First(&student)

	semesters := []models.Semester{
		{Year: "2021/2022", Term: "Ganjil", Active: false},
		{Year: "2021/2022", Term: "Genap", Active: false},
		{Year: "2022/2023", Term: "Ganjil", Active: false},
		{Year: "2022/2023", Term: "Genap", Active: false},
		{Year: "2023/2024", Term: "Ganjil", Active: false},
		{Year: "2023/2024", Term: "Genap", Active: false},
		{Year: "2024/2025", Term: "Ganjil", Active: true},
	}

	for _, s := range semesters {
		var exists models.Semester
		if err := db.Where("year = ? AND term = ?", s.Year, s.Term).First(&exists).Error; err == nil {
			log.Printf("Semester %s %s exists", s.Year, s.Term)
		} else {
			db.Create(&s)
			log.Printf("Semester %s %s created", s.Year, s.Term)
		}
	}

	var semesterList []models.Semester
	db.Find(&semesterList)

	krsData := map[string][]models.Course{
		"2021/2022-Ganjil": {
			{CourseCode: "IF101", CourseName: "Pengantar Informatika", SKS: 2, Lecturer: "Dr. Budi Santoso"},
			{CourseCode: "IF102", CourseName: "Matematika Dasar", SKS: 3, Lecturer: "Mila Rahma, M.Kom"},
		},
		"2021/2022-Genap": {
			{CourseCode: "IF103", CourseName: "Logika & Algoritma", SKS: 3, Lecturer: "Dr. Budi Santoso"},
			{CourseCode: "IF104", CourseName: "Pemrograman Dasar", SKS: 3, Lecturer: "Mila Rahma, M.Kom"},
		},
		// dst sesuai data
	}

	grades := []string{"A", "A-", "B+", "B", "B-", "C+", "C"}

	for _, sem := range semesterList {
		key := sem.Year + "-" + sem.Term

		var existingKRS models.KRS
		if err := db.Where("user_id = ? AND semester_id = ?", student.ID, sem.ID).First(&existingKRS).Error; err == nil {
			log.Printf("KRS %s exists, skipping...", key)
			continue
		}

		// Create KRS
		krs := models.KRS{UserID: student.ID, SemesterID: sem.ID, Finalized: true}
		db.Create(&krs)

		courses := krsData[key]

		// Buat Course & KRSDetail
		var krsDetails []models.KRSDetail
		for _, c := range courses {
			c.KRSID = krs.ID
			db.Create(&c)

			randomGrade := grades[rand.Intn(len(grades))]

			krsDetail := models.KRSDetail{
				KRSID:    krs.ID,
				CourseID: c.ID,
				Grade:    randomGrade,
			}
			db.Create(&krsDetail)
			krsDetails = append(krsDetails, krsDetail)
		}

		// Hitung GPA sederhana (misal A=4, A-=3.7, B+=3.3, B=3, dst)
		var totalSKS int
		var totalPoint float32
		gradeMap := map[string]float32{
			"A":  4.0,
			"A-": 3.7,
			"B+": 3.3,
			"B":  3.0,
			"B-": 2.7,
			"C+": 2.3,
			"C":  2.0,
		}
		for _, d := range krsDetails {
			totalSKS += d.Course.SKS
			totalPoint += float32(d.Course.SKS) * gradeMap[d.Grade]
		}
		gpa := totalPoint / float32(totalSKS)

		// Buat KHS dari semua KRSDetail
		// Buat KHS dari semua KRSDetail
		khs := models.KHS{
			UserID:     student.ID,
			SemesterID: sem.ID,
			GPA:        gpa,
		}
		db.Create(&khs)

		// === Tambahkan relasi ke KRSDetail === //
		for i := range krsDetails {
			krsDetails[i].KHSID = &khs.ID
			db.Model(&krsDetails[i]).Update("khs_id", khs.ID)
		}

		// Append details ke KHS
		db.Model(&khs).Association("Details").Append(krsDetails)

		log.Printf("KRS & KHS %s created with %d courses", key, len(courses))
	}

	log.Println("🎉 Seeder Selesai")
}
