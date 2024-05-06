package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type PayloadUser struct {
	Id    string   `json:"id"`
	Login string   `json:"login"`
	Teams []string `json:"teams"`
}

func GenerateJwtTokenString(id string, login string, teamIds []string, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"teams": teamIds,
	})
	return token.SignedString([]byte(key))
}

func ParseTokenString(tokenString, key string) (*PayloadUser, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &PayloadUser{
			Id:    claims["id"].(string),
			Login: claims["login"].(string),
			Teams: claims["teams"].([]string),
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
