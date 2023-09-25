package token

import (
	"github.com/stretchr/testify/require"
	"github.com/zsoltardai/simple_bank/util"
	"testing"
	"time"
)

func createPasetoMaker(t *testing.T) Maker {
	maker, err := NewPasetoMaker(util.RandomString(32))

	require.NoError(t, err)
	require.NotEmpty(t, maker)

	return maker
}

func TestPasetoMaker(t *testing.T) {
	maker := createPasetoMaker(t)

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

func TestExpiredPasetoToken(t *testing.T) {
	maker := createPasetoMaker(t)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
