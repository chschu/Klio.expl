package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"time"
)

type JwtGenerator interface {
	Generate(subject string, valid time.Duration) (jwtStr string, err error)
}

type JwtValidator interface {
	Validate(jwtStr string) (subject string, err error)
}

func NewJwtHandlers(method jwt.SigningMethod, generationKey any, validationKey any) (JwtGenerator, JwtValidator) {
	return &jwtGenerator{method: method, key: generationKey}, &jwtValidator{method: method, key: validationKey}
}

type jwtGenerator struct {
	method jwt.SigningMethod
	key    any
}

func (j *jwtGenerator) Generate(subject string, valid time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(valid)),
		Subject:   subject,
	}
	return jwt.NewWithClaims(j.method, claims).SignedString(j.key)

}

type jwtValidator struct {
	method jwt.SigningMethod
	key    any
}

func (j *jwtValidator) Validate(jwtStr string) (string, error) {
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
