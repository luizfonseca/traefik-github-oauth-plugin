package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTeams = []string{"team1", "team2"}

const (
	id    = "12345"
	login = "testuser"
	key   = "secretKey"
)

func TestGenerateJwtTokenString(t *testing.T) {
	// execution
	tokenString, err := GenerateJwtTokenString(id, login, testTeams, key, false)

	// assertion
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestParseTokenString(t *testing.T) {
	// setup
	tokenString, _ := GenerateJwtTokenString(id, login, testTeams, key, false)

	// execution
	payload, err := ParseTokenString(tokenString, key)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, id, payload.Id)
	assert.Equal(t, login, payload.Login)
	assert.Equal(t, testTeams, payload.Teams)
	assert.False(t, payload.TwoFactorEnabled)
}

func TestParseTokenString_EmptyTeams(t *testing.T) {
	// setup
	tokenString, _ := GenerateJwtTokenString(id, login, []string{}, key, false)

	// execution
	payload, err := ParseTokenString(tokenString, key)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, id, payload.Id)
	assert.Equal(t, login, payload.Login)
	assert.Equal(t, payload.Teams, []string{})
	assert.False(t, payload.TwoFactorEnabled)
}

func TestParseTokenString_NoTeams(t *testing.T) {
	// setup
	tokenString, _ := GenerateJwtTokenString(id, login, nil, key, false)

	// execution
	payload, err := ParseTokenString(tokenString, key)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, id, payload.Id)
	assert.Equal(t, login, payload.Login)
	assert.Equal(t, payload.Teams, []string{})
	assert.False(t, payload.TwoFactorEnabled)
}

func TestParseTokenString_With2FAEnabled(t *testing.T) {
	// setup
	tokenString, _ := GenerateJwtTokenString(id, login, nil, key, true)

	// execution
	payload, err := ParseTokenString(tokenString, key)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, id, payload.Id)
	assert.Equal(t, login, payload.Login)
	assert.Equal(t, payload.Teams, []string{})
	assert.True(t, payload.TwoFactorEnabled)
}

func TestParseTokenString_InvalidToken(t *testing.T) {
	// setup
	tokenString := "invalidtoken"

	// execution
	payload, err := ParseTokenString(tokenString, key)

	// assertion
	assert.Error(t, err)
	assert.Nil(t, payload)
}

func TestParseTokenString_InvalidKey(t *testing.T) {
	// setup
	tokenString, _ := GenerateJwtTokenString(id, login, testTeams, key, false)
	invalidKey := "invalidkey"

	// execution
	payload, err := ParseTokenString(tokenString, invalidKey)

	// assertion
	assert.Error(t, err)
	assert.Nil(t, payload)
}
