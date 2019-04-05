package auth

import "context"

type key int

const (
	secretKey key = 0
	userKey   key = 1
)

func NewContextWithSecretKey(ctx context.Context, secret string) context.Context {
	return context.WithValue(ctx, secretKey, secret)
}

func SecretKeyFromContext(ctx context.Context) (string, bool) {
	secret, ok := ctx.Value(secretKey).(string)
	return secret, ok
}

func NewContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userKey).(string)
	return userID, ok
}
