package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtGenerate func(subject string, valid time.Duration) (jwtStr string, err error)
type JwtValidate func(jwtStr string) (subject string, err error)

func NewJwtHandlers() (JwtGenerate, JwtValidate, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return jwtGenerateFunc(private), jwtValidateFunc(public), nil
}

func jwtGenerateFunc(k ed25519.PrivateKey) JwtGenerate {
	return func(subject string, valid time.Duration) (string, error) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(valid)),
			Subject:   subject,
		}
		return jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(k)
	}
}

func jwtValidateFunc(k ed25519.PublicKey) JwtValidate {
	return func(jwtStr string) (string, error) {
		token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return k, nil
		})
		if err != nil {
			return "", err
		}
		if !token.Valid {
			return "", errors.New("token is invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", errors.New("token does not contain map claims")
		}
		subject, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("sub claim is not a string")
		}
		return subject, nil
	}
}
