package datastore

import (
	"context"
	"crypto/rand"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(ctx context.Context, password string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "hashPassword")
	defer span.Finish()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func encodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
