package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type PayloadUser struct {
	Id    string   `json:"id"`
	Login string   `json:"login"`
	Teams []string `json:"teams"`
}

func GenerateJwtTokenString(id string, login string, teamIds []string, key string, exp time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"teams": teamIds,
		// buffer of time to expire token is 10 seconds from the set time
		"exp": jwt.NewNumericDate(exp.Add(time.Second * 60)),
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
		var teamFromClaims []interface{}

		// Check for expiration time
		// if claims.Valid() != nil {
		// 	return nil, fmt.Errorf("token is expired")
		// }

		switch claims["teams"].(type) {
		case []interface{}:
			teamFromClaims = claims["teams"].([]interface{})
		case nil:
			// Backwards compatible with previous tokens, nothing to do
		}

		teams := make([]string, len(teamFromClaims))

		for i, v := range teamFromClaims {
			if stringValue, ok := v.(string); ok {
				teams[i] = stringValue
			}
		}

		return &PayloadUser{
			Id:    claims["id"].(string),
			Login: claims["login"].(string),
			Teams: teams,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
