package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject:   userID.String(),
		})
	s, err := t.SignedString([]byte(tokenSecret))
	return s, err
}

func ValidateJwt(tokenString, tokenSecrect string) (uuid.UUID, error) {
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecrect), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("Authorization header missing")
	}
	const prefix = "Bearer"
	if !strings.HasPrefix(authorizationHeader[0], prefix) {
		return "", fmt.Errorf("Authorization header format must be 'Bearer TOKEN' ")
	}
	token := strings.Split(authorizationHeader[0], " ")[1]
	return token, nil
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authorizationHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("Authorization header missing")
	}
	const prefix = "ApiKey"
	if !strings.HasPrefix(authorizationHeader[0], prefix) {
		return "", fmt.Errorf("Authorization header format must be 'ApiKey' TOKEN")
	}
	token := strings.Split(authorizationHeader[0], " ")[1]
	return token, nil
}
