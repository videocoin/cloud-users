package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/dchest/authcookie"
	"github.com/opentracing/opentracing-go"
	ds "github.com/videocoin/cloud-users/datastore"
	"golang.org/x/crypto/bcrypt"
)

var MinTokenLength = authcookie.MinLength

var (
	ErrMalformedToken = errors.New("malformed token")
	ErrExpiredToken   = errors.New("token expired")
	ErrWrongSignature = errors.New("wrong token signature")
)

func checkPasswordHash(ctx context.Context, password, hash string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "checkPasswordHash")
	defer span.Finish()

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	span.LogKV("is equal", err == nil)

	return err == nil
}

func getUserSecretKey(passwordHash, secret []byte) []byte {
	m := hmac.New(sha256.New, secret)
	_, _ = m.Write(passwordHash)
	return m.Sum(nil)
}

func getSignature(b []byte, secret []byte) []byte {
	keym := hmac.New(sha256.New, secret)
	_, _ = keym.Write(b)
	m := hmac.New(sha256.New, keym.Sum(nil))
	_, _ = m.Write(b)
	return m.Sum(nil)
}

// NewToken returns a new password reset token for the given login, which
// expires after the given time duration since now, signed by the key generated
// from the given password value (which can be any value that will be changed
// once a user resets their password, such as password hash or salt used to
// generate it), and the given secret key.
func newRecoveryToken(email string, duration time.Duration, passwordHash, secret []byte) string {
	secretKey := getUserSecretKey(passwordHash, secret)
	return authcookie.NewSinceNow(email, duration, secretKey)
}

// VerifyToken verifies the given token with the password value returned by the
// given function and the given secret key, and returns login extracted from
// the valid token. If the token is not valid, the function returns an error.
//
// Function pwdvalFn must return the current password value for the login it
// receives in arguments, or an error. If it returns an error, VerifyToken
// returns the same error.
func verifyRecoveryToken(ctx context.Context, token string, getUserFunc func(context.Context, string) (*ds.User, error), secret []byte) (user *ds.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "verifyRecoveryToken")
	defer span.Finish()

	blen := base64.URLEncoding.DecodedLen(len(token))
	// Avoid allocation if the token is too short
	if blen <= 4+32 {
		err = ErrMalformedToken
		return
	}
	b := make([]byte, blen)
	blen, err = base64.URLEncoding.Decode(b, []byte(token))
	if err != nil {
		return
	}
	// Decoded length may be bifferent from max length, which
	// we allocated, so check it, and set new length for b
	if blen <= 4+32 {
		err = ErrMalformedToken
		return
	}
	b = b[:blen]

	data := b[:blen-32]
	exp := time.Unix(int64(binary.BigEndian.Uint32(data[:4])), 0)
	if exp.Before(time.Now()) {
		err = ErrExpiredToken
		return
	}

	email := string(data[4:])

	span.LogKV("email", email)

	user, err = getUserFunc(ctx, email)
	if err != nil {
		return
	}

	sig := b[blen-32:]
	sk := getUserSecretKey([]byte(user.Password), secret)
	realSig := getSignature(data, sk)
	if subtle.ConstantTimeCompare(realSig, sig) != 1 {
		err = ErrWrongSignature
		return
	}

	return
}
