package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/laurawarren88/go_spa_backend.git/models"
)

type Claims struct {
	UserID   uint   `json:"sub"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	accessSecret := os.Getenv("ACCESS_SECRET_KEY")
	if accessSecret == "" {
		log.Println("ACCESS_SECRET_KEY is not set")
		return "", fmt.Errorf("access secret key not set in environment")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}

func GenerateRefreshToken(user models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}

	refreshSecret := os.Getenv("REFRESH_SECRET_KEY")
	if refreshSecret == "" {
		log.Println("REFRESH_SECRET_KEY is not set")
		return "", fmt.Errorf("refresh secret key not set in environment")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}
