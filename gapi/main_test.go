package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/token"
	"github.com/Hypersus/simplebank/util"
	"github.com/Hypersus/simplebank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenKey:      util.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.TokenMaker, username string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.GenerateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
