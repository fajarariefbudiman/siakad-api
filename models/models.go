package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nim       string         `gorm:"type:varchar(50);uniqueIndex" json:"nim"`
	Password  string         `gorm:"type:varchar(255)" json:"-"`
	Role      string         `gorm:"type:varchar(20)" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Semester struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Year      string    `gorm:"type:varchar(20)" json:"year"`
	Term      string    `gorm:"type:varchar(20)" json:"term"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KRS → menyimpan daftar mata kuliah yang diambil mahasiswa
type KRS struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	UserID     uint        `json:"user_id"`
	SemesterID uint        `json:"semester_id"`
	Finalized  bool        `json:"finalized"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Details    []KRSDetail `json:"details" gorm:"foreignKey:KRSID"` // KRSDetail berisi grade dan nilai
}

// Course → hanya untuk input, tidak menyimpan grade
type Course struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	KRSID      uint      `json:"krs_id"` // relasi ke KRS
	CourseCode string    `gorm:"type:varchar(20)" json:"course_code"`
	CourseName string    `gorm:"type:varchar(150)" json:"course_name"`
	Lecturer   string    `gorm:"type:varchar(100)" json:"lecturer"`
	SKS        int       `json:"sks"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// KRSDetail → menyimpan nilai per mata kuliah dari KRS
type KRSDetail struct {
	ID    uint  `gorm:"primaryKey" json:"id"`
	KRSID uint  `json:"krs_id"`
	KHSID *uint `json:"khs_id"` // tambahkan ini

	CourseID uint   `json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID"`

	Grade     string    `json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KHS → merujuk ke KRSDetail untuk menampilkan nilai
type KHS struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	UserID     uint        `json:"user_id"`
	SemesterID uint        `json:"semester_id"`
	GPA        float32     `json:"gpa"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Details    []KRSDetail `json:"details" gorm:"foreignKey:KHSID"`
}

// Payment dan Post tetap sama
type Payment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Paid        bool      `json:"paid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Post struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"type:varchar(150)" json:"title"`
	Slug      string    `gorm:"type:varchar(150);uniqueIndex" json:"slug"`
	Body      string    `gorm:"type:text" json:"body"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TugasAkhir struct {
	ID uint `gorm:"primarykey" json:"id"`

	// Relasi mahasiswa
	UserID uint `json:"user_id"`
	User   User `json:"user"`

	// Jenis: skripsi / kp
	Category string `gorm:"type:varchar(20)" json:"category"`

	// Informasi mahasiswa
	StudentName string `gorm:"type:varchar(100)" json:"student_name"`
	StudentNim  string `gorm:"type:varchar(50)" json:"student_nim"`

	// Informasi penelitian
	Title         string `gorm:"type:varchar(255)" json:"title"`
	ResearchPlace string `gorm:"type:varchar(255)" json:"research_place"`
	Unit          string `gorm:"type:varchar(100)" json:"unit"`
	Address       string `gorm:"type:varchar(255)" json:"address"`

	// Kontak perusahaan
	CompanyEmail string `gorm:"type:varchar(100)" json:"company_email"`
	CompanyPhone string `gorm:"type:varchar(30)" json:"company_phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
