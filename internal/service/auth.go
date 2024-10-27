package service

import (
	"os"

	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/repository"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const(
	tokenTTL = 12 * time.Hour
)

var(
	signingKey = os.Getenv("SIGNING_KEY")
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct{
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService{
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user recipes.User) (int, error){
	var err error

	user.Salt, err = GenerateSalt()
	if err != nil{
		return 0, err
	}

	user.Password = HashPassword(user.Password, user.Salt)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	salt, err := a.repo.GetUserSalt(username)

	if err != nil{
		return "User does not exist", err
	} else {
		user, err := a.repo.GetUser(username, HashPassword(password, salt))
		if err != nil {
			return "", err
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt: time.Now().Unix(),
			},
			user.Id,
		})
		
		return token.SignedString([]byte(signingKey))
	}
	
}

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword - функция для хеширования пароля с использованием соли
func HashPassword(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}


