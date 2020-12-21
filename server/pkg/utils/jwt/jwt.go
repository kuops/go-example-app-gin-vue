package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var  jwtSignKey = []byte("qwerasdf")
var  TokenExpired = errors.New("token expired")

type CustomClaims struct {
	ID uint64
	UUID string
	Name string
	jwt.StandardClaims
}

func CreateToken(cliams *CustomClaims) (string, error) {
	cliams.Issuer = "go-example-app"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	tokenString,err := token.SignedString(jwtSignKey)
	if err != nil {
		return "", err
	}
	return tokenString,nil
}

func ParseToken(t string) (*CustomClaims,error) {
	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSignKey, nil
	})

	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return nil, TokenExpired
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil,err
}