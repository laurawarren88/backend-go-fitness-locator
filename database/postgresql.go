package database

import (
	"fmt"
	"log"
	"os"

	"github.com/laurawarren88/go_spa_backend.git/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	// Change db to localhost if running locally of machine and not from docker
	// dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
	// dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
	log.Printf("DSN: %s", dsn)

	if err := DB.AutoMigrate(&models.User{}, &models.Place{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	// log.Println("Database migration completed")
}

func GetDB() *gorm.DB {
	return DB
}

func SetupAdminUser(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Password hashing failed: %v", err)
		return err
	}
	// log.Println("Hashed Password:", string(hashedPassword))

	admin := models.User{
		Username: "admin",
		Email:    "admin@admin.com",
		Password: string(hashedPassword),
		IsAdmin:  true,
	}

	var existingUser models.User
	result := db.Where("email = ?", admin.Email).First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	} else if result.Error == nil {
		existingUser.Username = admin.Username
		existingUser.Password = admin.Password
		existingUser.IsAdmin = admin.IsAdmin
		if err := db.Save(&existingUser).Error; err != nil {
			return err
		}
		log.Println("Admin user updated successfully")
	} else {
		return result.Error
	}

	return nil
}
