package token

import (
	"testing"
	"time"

	"github.com/Hypersus/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestjwtMaker(t *testing.T) {
	// generate a random secret key
	var secretKey = util.RandomString(64)
	// construct a new jwtMaker
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	// generate a random username
	var username = util.RandomOwner()
	// generate a random duration
	var duration = time.Minute

	var IssuedAt = time.Now()
	var ExpiredAt = IssuedAt.Add(duration)
	// generate a token for the username
	token, err := jwtMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// validate the token
	payload, err := jwtMaker.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	// generate a random secret key
	var secretKey = util.RandomString(64)
	// construct a new jwtMaker
	jwtMaker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	// generate a random username
	var username = util.RandomOwner()
	// generate a random duration
	var duration = -time.Minute

	// generate a token for the username
	token, err := jwtMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// validate the token
	payload, err := jwtMaker.ValidateToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
	require.Equal(t, ErrTokenExpired, err)
}
