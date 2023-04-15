package token

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrTokenExpired      = errors.New("token is expired")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUnregisteredMaker = errors.New("unregistered token maker")
	ErrUnsafekey         error
)

func UnsafeKey(err error) error {
	ErrUnsafekey = fmt.Errorf("unsafe key: %v", err)
	return ErrUnsafekey
}

var MakerFactory = make(map[string]func(secretKey string) (TokenMaker, error))
var signingMethodLock = new(sync.RWMutex)

func init() {
	RegisterMaker("JWT", NewJWTMaker)
}

func RegisterMaker(name string, f func(secretKey string) (TokenMaker, error)) {
	signingMethodLock.Lock()
	defer signingMethodLock.Unlock()

	MakerFactory[name] = f
}

func NewMaker(name string) (f func(secretKey string) (TokenMaker, error)) {
	signingMethodLock.RLock()
	defer signingMethodLock.RUnlock()

	if f, ok := MakerFactory[name]; ok {
		return f
	}
	return nil
}

type TokenMaker interface {
	// GenerateToken generates a new token for the given user.
	GenerateToken(username string, duration time.Duration) (string, error)
	// ValidateToken validates the given token and returns the payload
	// associated with the token if failed return error.
	ValidateToken(token string) (*Payload, error)
}
