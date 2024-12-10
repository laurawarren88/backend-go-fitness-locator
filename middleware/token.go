package middleware

// import (
// 	"os"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/laurawarren88/go_spa_backend.git/models"
// )

// func generateToken(user models.User) (string, error) {
// 	claims := Claims{
// 		UserID:   user.ID,
// 		Username: user.Username,
// 		IsAdmin:  user.IsAdmin,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
// }
