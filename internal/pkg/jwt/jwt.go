package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type PayloadUser struct {
	Id               string   `json:"id"`
	Login            string   `json:"login"`
	Teams            []string `json:"teams"`
	TwoFactorEnabled bool     `json:"two_factor_enabled"`
}

func GenerateJwtTokenString(id string, login string, teamIds []string, key string, two_factor_enabled bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                 id,
		"login":              login,
		"teams":              teamIds,
		"two_factor_enabled": two_factor_enabled,
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

		twoFactorEnabled := false
		if claims["two_factor_enabled"] != nil {
			if factorEnabled, ok := claims["two_factor_enabled"].(bool); ok {
				twoFactorEnabled = factorEnabled
			}
		}

		return &PayloadUser{
			Id:               claims["id"].(string),
			Login:            claims["login"].(string),
			Teams:            teams,
			TwoFactorEnabled: twoFactorEnabled,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
