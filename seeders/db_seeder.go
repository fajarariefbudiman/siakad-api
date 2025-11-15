package seeders

import (
	"api-siakad/config"
	"api-siakad/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	db := config.DB
	log.Println("Running Seeder...")
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := []models.User{
		{
			Name:     "Rahman Fajar Banyu Adji",
			Email:    "rahman@student.com",
			Nim:      "230441001",
			Password: string(hash),
			Role:     "student",
		},
		{
			Name:     "Dr. Budi Santoso",
			Email:    "budi@dosen.com",
			Nim:      "DSN001",
			Password: string(hash),
			Role:     "dosen",
		},
		{
			Name:     "Mila Rahma, M.Kom",
			Email:    "mila@dosen.com",
			Nim:      "DSN002",
			Password: string(hash),
			Role:     "dosen",
		},
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

	// Ambil user mahasiswa
	var student models.User
	db.Where("email = ?", "rahman@student.com").First(&student)

	semesters := []models.Semester{
		{Year: "2021/2022", Term: "Ganjil", Active: false}, // S1
		{Year: "2021/2022", Term: "Genap", Active: false},  // S2
		{Year: "2022/2023", Term: "Ganjil", Active: false}, // S3
		{Year: "2022/2023", Term: "Genap", Active: false},  // S4
		{Year: "2023/2024", Term: "Ganjil", Active: false}, // S5
		{Year: "2023/2024", Term: "Genap", Active: false},  // S6
		{Year: "2024/2025", Term: "Ganjil", Active: true},  // S7 (aktif)
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

	// Ambil semester aktif (S7)
	var activeSemester models.Semester
	db.Where("active = ?", true).First(&activeSemester)

	// ============================
	// 3. KRS + DETAIL KRS per semester
	// ============================

	krsData := map[string][]models.KRSDetail{
		"2021/2022-Ganjil": {
			{CourseCode: "IF101", CourseName: "Pengantar Informatika", SKS: 2},
			{CourseCode: "IF102", CourseName: "Matematika Dasar", SKS: 3},
		},
		"2021/2022-Genap": {
			{CourseCode: "IF103", CourseName: "Logika & Algoritma", SKS: 3},
			{CourseCode: "IF104", CourseName: "Pemrograman Dasar", SKS: 3},
		},
		"2022/2023-Ganjil": {
			{CourseCode: "IF201", CourseName: "Struktur Data", SKS: 3},
			{CourseCode: "IF202", CourseName: "Basis Data", SKS: 3},
		},
		"2022/2023-Genap": {
			{CourseCode: "IF203", CourseName: "Sistem Informasi", SKS: 2},
			{CourseCode: "IF204", CourseName: "Analisis Sistem", SKS: 3},
		},
		"2023/2024-Ganjil": {
			{CourseCode: "IF301", CourseName: "Pemrograman Web", SKS: 3},
			{CourseCode: "IF302", CourseName: "Jaringan Komputer", SKS: 3},
		},
		"2023/2024-Genap": {
			{CourseCode: "IF303", CourseName: "Manajemen Proyek IT", SKS: 3},
		},
		"2024/2025-Ganjil": {
			{CourseCode: "IF401", CourseName: "Mobile Programming", SKS: 3},
			{CourseCode: "IF402", CourseName: "Machine Learning", SKS: 3},
		},
	}

	var semesterList []models.Semester
	db.Find(&semesterList)

	for _, sem := range semesterList {
		key := sem.Year + "-" + sem.Term

		// cek krs
		var existingKRS models.KRS
		if err := db.Where("user_id = ? AND semester_id = ?", student.ID, sem.ID).First(&existingKRS).Error; err == nil {
			log.Printf("KRS %s exists, skipping...", key)
			continue
		}

		krs := models.KRS{
			UserID:     student.ID,
			SemesterID: sem.ID,
			Finalized:  true,
		}
		db.Create(&krs)

		// tambah details
		for _, d := range krsData[key] {
			d.KRSID = krs.ID
			db.Create(&d)
		}

		log.Printf("KRS %s created", key)
	}

	// ============================
	// 4. KHS per semester
	// ============================
	khsData := []models.KHS{
		{UserID: student.ID, SemesterID: semesterList[0].ID, GPA: 3.10},
		{UserID: student.ID, SemesterID: semesterList[1].ID, GPA: 3.20},
		{UserID: student.ID, SemesterID: semesterList[2].ID, GPA: 3.25},
		{UserID: student.ID, SemesterID: semesterList[3].ID, GPA: 3.30},
		{UserID: student.ID, SemesterID: semesterList[4].ID, GPA: 3.40},
		{UserID: student.ID, SemesterID: semesterList[5].ID, GPA: 3.45},
		{UserID: student.ID, SemesterID: activeSemester.ID, GPA: 3.50},
	}

	for _, khs := range khsData {
		var exists models.KHS
		if err := db.Where("user_id = ? AND semester_id = ?", khs.UserID, khs.SemesterID).First(&exists).Error; err == nil {
			continue
		}
		db.Create(&khs)
	}

	// ============================
	// 5. PAYMENTS
	// ============================
	payments := []models.Payment{
		{UserID: student.ID, Amount: 2500000, Description: "Pembayaran Semester 1", Paid: true},
		{UserID: student.ID, Amount: 2500000, Description: "Pembayaran Semester 2", Paid: true},
		{UserID: student.ID, Amount: 2750000, Description: "Pembayaran Semester 3", Paid: true},
		{UserID: student.ID, Amount: 2750000, Description: "Pembayaran Semester 4", Paid: true},
		{UserID: student.ID, Amount: 3000000, Description: "Pembayaran Semester 5", Paid: false},
		{UserID: student.ID, Amount: 3000000, Description: "Pembayaran Semester 6", Paid: false},
	}

	for _, p := range payments {
		var exists models.Payment
		if err := db.Where("user_id = ? AND description = ?", p.UserID, p.Description).First(&exists).Error; err == nil {
			continue
		}
		db.Create(&p)
	}

	// ============================
	// 6. POSTS
	// ============================
	posts := []models.Post{
		{
			Title:     "Pengumuman Libur Nasional",
			Slug:      "libur-nasional",
			Body:      "Kampus akan libur pada tanggal 17 Agustus.",
			Published: true,
		},
		{
			Title:     "Pendaftaran Wisuda 2025",
			Slug:      "pendaftaran-wisuda",
			Body:      "Pendaftaran wisuda telah dibuka hingga 30 Juni.",
			Published: true,
		},
	}

	for _, post := range posts {
		var exists models.Post
		if err := db.Where("slug = ?", post.Slug).First(&exists).Error; err == nil {
			continue
		}
		db.Create(&post)
	}

	log.Println("🎉 Seeder Selesai")
}
