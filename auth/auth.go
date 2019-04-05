package auth

import (
	"context"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var (
	ErrInvalidSecret = errors.New("invalid secret")
	ErrInvalidToken  = errors.New("invalid token")
)

type JWTClaims struct {
	UserID string `json:"aid"`
	jwt.StandardClaims
}

func AuthFromContext(ctx context.Context) (context.Context, error) {
	secret, ok := SecretKeyFromContext(ctx)
	if !ok {
		return ctx, ErrInvalidSecret
	}

	jwtToken, err := grpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, ErrInvalidToken
	}

	t, err := jwt.ParseWithClaims(jwtToken, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return ctx, err
	}

	if !t.Valid {
		return ctx, ErrInvalidToken
	}

	ctx = NewContextWithUserID(ctx, t.Claims.(*JWTClaims).UserID)

	return ctx, nil
}
