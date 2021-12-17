package msAPI

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type MsJWT interface {
	GenerateToken(claims jwt.Claims) string
	GetClaims(claims jwt.Claims, tokenStr string) error
	VerifyToken(tokenStr string) (*jwt.Token, error)
}
type msJWT struct {
	secret string
}

func NewJWT(secret string) MsJWT {
	return &msJWT{
		secret: secret,
	}
}

func GetJWTDefaultSecret() string {
	secret := os.Getenv("jwt-secretKey")
	if secret == "" {
		secret = "jwt-secretKey"
	}
	return secret
}

func (j *msJWT) GenerateToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		panic(err)
	}

	return tokenStr
}

func (j *msJWT) GetClaims(claims jwt.Claims, tokenStr string) error {
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	return err
}

func (j *msJWT) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
}
