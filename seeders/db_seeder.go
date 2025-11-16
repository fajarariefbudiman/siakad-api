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

	// === USERS === //
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

	// === SEMESTERS === //
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

	// === COURSE DATA === //
	krsData := map[string][]struct {
		CourseCode string
		CourseName string
		Lecturer   string
		SKS        int
	}{
		"2021/2022-Ganjil": {
			{"IF101", "Pengantar Teknologi Informasi", "Dr. Andi", 3},
			{"IF102", "Algoritma & Pemrograman I", "Dr. Budi", 4},
			{"IF103", "Matematika Dasar", "Dr. Wati", 3},
			{"IF104", "Bahasa Indonesia", "Dr. Rina", 2},
		},
		"2021/2022-Genap": {
			{"IF201", "Struktur Data", "Dr. Budi", 4},
			{"IF202", "Basis Data I", "Dr. Andi", 3},
			{"IF203", "Matematika Diskrit", "Dr. Dewi", 3},
			{"IF204", "Pancasila", "Dr. Rina", 2},
		},
		"2022/2023-Ganjil": {
			{"IF301", "Basis Data II", "Dr. Andi", 3},
			{"IF302", "Pemrograman Berorientasi Objek", "Dr. Budi", 4},
			{"IF303", "Jaringan Komputer", "Dr. Wawan", 3},
			{"IF304", "Statistika", "Dr. Sari", 2},
		},
		"2022/2023-Genap": {
			{"IF401", "Sistem Operasi", "Dr. Hasan", 3},
			{"IF402", "Pemrograman Web", "Dr. Budi", 4},
			{"IF403", "Rekayasa Perangkat Lunak", "Dr. Andi", 3},
			{"IF404", "Kewirausahaan", "Dr. Tono", 2},
		},
		"2023/2024-Ganjil": {
			{"IF501", "Machine Learning", "Dr. Sari", 3},
			{"IF502", "Pemrograman Mobile", "Dr. Budi", 4},
			{"IF503", "Keamanan Informasi", "Dr. Wawan", 3},
			{"IF504", "Manajemen Proyek TI", "Dr. Andi", 2},
		},
		"2023/2024-Genap": {
			{"IF601", "Cloud Computing", "Dr. Hasan", 3},
			{"IF602", "Data Mining", "Dr. Sari", 4},
			{"IF603", "UI / UX Design", "Dr. Rina", 3},
			{"IF604", "Etika Profesi", "Dr. Tono", 2},
		},
		"2024/2025-Ganjil": {
			{"IF701", "Kecerdasan Buatan Lanjut", "Dr. Hasan", 3},
			{"IF702", "Fullstack Development", "Dr. Budi", 4},
			{"IF703", "Data Visualization", "Dr. Sari", 3},
			{"IF704", "Metode Penelitian", "Dr. Andi", 2},
		},
	}

	// === NILAI RANDOM === //
	grades := []string{"A", "A-", "B+", "B", "B-", "C+", "C"}

	// === LOOP SEMESTER === //
	for _, sem := range semesterList {
		key := sem.Year + "-" + sem.Term

		// Cek existing KRS
		var existingKRS models.KRS
		if err := db.Where("user_id = ? AND semester_id = ?", student.ID, sem.ID).First(&existingKRS).Error; err == nil {
			log.Printf("KRS %s exists, skipping...", key)
			continue
		}

		// Buat KRS
		krs := models.KRS{UserID: student.ID, SemesterID: sem.ID, Finalized: true}
		db.Create(&krs)

		courses := krsData[key]
		if len(courses) == 0 {
			log.Printf("SKIP %s — tidak ada mata kuliah", key)
			continue
		}

		// Buat Course & KRSDetail
		var details []models.KRSDetail
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

		for _, c := range courses {
			// Insert Course
			course := models.Course{
				KRSID:      krs.ID,
				CourseCode: c.CourseCode,
				CourseName: c.CourseName,
				Lecturer:   c.Lecturer,
				SKS:        c.SKS,
			}
			db.Create(&course)

			// Nilai random
			grade := grades[rand.Intn(len(grades))]

			// Insert KRSDetail
			krsDetail := models.KRSDetail{
				KRSID:    krs.ID,
				CourseID: course.ID,
				Grade:    grade,
			}
			db.Create(&krsDetail)

			details = append(details, krsDetail)

			// Hitung GPA
			totalSKS += course.SKS
			totalPoint += float32(course.SKS) * gradeMap[grade]
		}

		// GPA final
		gpa := totalPoint / float32(totalSKS)

		// Buat KHS
		khs := models.KHS{
			UserID:     student.ID,
			SemesterID: sem.ID,
			GPA:        gpa,
		}
		db.Create(&khs)

		log.Printf("KRS & KHS %s created with %d courses", key, len(courses))
	}

	paymentData := map[string]float64{
		"2021/2022-Ganjil": 10000000,
		"2021/2022-Genap":  8000000,
		"2022/2023-Ganjil": 7850000,
		"2022/2023-Genap":  8200000,
		"2023/2024-Ganjil": 9000000,
		"2023/2024-Genap":  9500000,
		"2024/2025-Ganjil": 10000000,
	}

	for _, sem := range semesterList {
		key := sem.Year + "-" + sem.Term

		total, exists := paymentData[key]
		if !exists {
			continue
		}

		// Check existing payments
		var existing []models.Payment
		if err := db.Where("user_id = ? AND semester_id = ?", student.ID, sem.ID).Find(&existing).Error; err == nil && len(existing) > 0 {
			log.Printf("Payments for %s already exist, skipping...", key)
			continue
		}

		// Buat payment, misal 50% terbayar
		paidAmount := total * 0.5

		payment := models.Payment{
			UserID:      student.ID,
			SemesterID:  sem.ID,
			Amount:      total,
			Paid:        true, // bisa disesuaikan
			Description: "Pembayaran " + key,
		}
		db.Create(&payment)

		log.Printf("Payment for %s created: total %.0f, paid %.0f", key, total, paidAmount)
	}
	log.Println("🎉 Seeder Selesai")
}
