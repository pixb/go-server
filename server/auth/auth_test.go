package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pixb/go-server/store"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	// Test data
	userID := int64(1)
	username := "testuser"
	role := "user"
	secret := "testsecret"

	// Generate access token
	token, err := GenerateAccessToken(userID, username, role, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate access token
	claims, err := ValidateAccessToken(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.NotEmpty(t, claims.ExpiresAt)
	assert.NotEmpty(t, claims.IssuedAt)
	assert.NotEmpty(t, claims.NotBefore)
	assert.Equal(t, "go-server", claims.Issuer)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	// Test data
	invalidToken := "invalidtoken"
	secret := "testsecret"

	// Validate invalid access token
	claims, err := ValidateAccessToken(invalidToken, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	// Test data
	userID := int64(1)
	username := "testuser"
	role := "user"
	secret := "testsecret"

	// Generate access token with short expiration
	oldTime := time.Now().Add(-24 * time.Hour)
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(oldTime),
			IssuedAt:  jwt.NewNumericDate(oldTime),
			NotBefore: jwt.NewNumericDate(oldTime),
			Issuer:    "go-server",
		},
		UserID:   userID,
		Username: username,
		Role:     role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Validate expired access token
	validatedClaims, err := ValidateAccessToken(tokenString, secret)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
}

func TestHashPassword(t *testing.T) {
	// Test data
	password := "testpassword"

	// Hash password
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Check password
	assert.True(t, CheckPassword(password, hashedPassword))
	assert.False(t, CheckPassword("wrongpassword", hashedPassword))
}

func TestAuthenticator_Authenticate(t *testing.T) {
	// Test data
	userID := int64(1)
	username := "testuser"
	role := "user"
	secret := "testsecret"

	// Generate access token
	token, err := GenerateAccessToken(userID, username, role, secret)
	assert.NoError(t, err)

	// Create authenticator
	authenticator := NewAuthenticator(&store.Store{}, secret)

	// Test authentication with valid token
	result := authenticator.Authenticate(nil, "Bearer "+token)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Claims)
	assert.Equal(t, userID, result.Claims.UserID)
	assert.Equal(t, username, result.Claims.Username)
	assert.Equal(t, role, result.Claims.Role)

	// Test authentication with invalid token
	result = authenticator.Authenticate(nil, "Bearer invalidtoken")
	assert.Nil(t, result)

	// Test authentication with empty header
	result = authenticator.Authenticate(nil, "")
	assert.Nil(t, result)
}
