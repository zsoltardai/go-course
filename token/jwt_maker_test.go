package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/zsoltardai/simple_bank/util"
	"testing"
	"time"
)

func createJWTMaker(t *testing.T) Maker {
	maker, err := NewJWTMaker(util.RandomString(32))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	return maker
}

func TestJWTMaker(t *testing.T) {
	maker := createJWTMaker(t)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker := createJWTMaker(t)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	maker := createJWTMaker(t)

	payload, err := NewPayload(util.RandomOwner(), time.Minute)

	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
