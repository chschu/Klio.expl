package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"time"
)

func NewJWTGenerator(method jwt.SigningMethod, generationKey any) *jwtGenerator {
	return &jwtGenerator{method: method, key: generationKey}
}

func NewJWTValidator(method jwt.SigningMethod, validationKey any) *jwtValidator {
	return &jwtValidator{method: method, key: validationKey}
}

type jwtGenerator struct {
	method jwt.SigningMethod
	key    any
}

func (j *jwtGenerator) GenerateJWT(subject string, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   subject,
	}
	return jwt.NewWithClaims(j.method, claims).SignedString(j.key)

}

type jwtValidator struct {
	method jwt.SigningMethod
	key    any
}

func (j *jwtValidator) ValidateJWT(jwtStr string) (string, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != j.method.Alg() {
			return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Method.Alg())
		}
		return j.key, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token is invalid")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims type: %v", reflect.TypeOf(token.Claims))
	}
	return claims.Subject, nil
}
