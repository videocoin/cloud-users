package datastore

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(ctx context.Context, password string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "hashPassword")
	defer span.Finish()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
