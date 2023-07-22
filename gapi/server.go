package gapi

import (
	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/pb"
	"github.com/Hypersus/simplebank/token"
	"github.com/Hypersus/simplebank/util"
	"github.com/Hypersus/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	// data layer
	store db.Store
	// token maker for authentication
	tokenMaker token.TokenMaker
	// configuration
	config          util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, distributor worker.TaskDistributor) (*Server, error) {
	Maker := token.NewMaker(config.TokenType)
	if Maker == nil {
		return nil, token.ErrUnregisteredMaker
	}
	tokenMaker, err := Maker(config.TokenKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: distributor,
	}
	return server, nil
}
