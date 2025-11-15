package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseIntEnv(key string, fallback int) int {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Invalid int env %s: %v", key, err)
		return fallback
	}
	return i
}

func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "secret")
}

func GetJWTExpiry() time.Duration {
	h := ParseIntEnv("JWT_EXPIRE_HOURS", 72)
	return time.Hour * time.Duration(h)
}
