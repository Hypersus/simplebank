package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Hypersus/simplebank/api"
	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/gapi"
	"github.com/Hypersus/simplebank/mail"
	"github.com/Hypersus/simplebank/pb"
	"github.com/Hypersus/simplebank/util"
	"github.com/Hypersus/simplebank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DevMode     = "dev"
	ReleaseMode = "release"
	Local       = "127.0.0.1"
	Global      = "0.0.0.0"
)

func setMode(c *util.Config) {
	switch c.Mode {
	case DevMode:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		c.HTTPServerAddress = Local + c.HTTPServerAddress
		c.GRPCServerAddress = Local + c.GRPCServerAddress
		c.RedisAddress = Local + c.RedisAddress
	default:
		c.HTTPServerAddress = Global + c.HTTPServerAddress
		c.GRPCServerAddress = Global + c.GRPCServerAddress
		c.RedisAddress = Global + c.RedisAddress
	}
}

func main() {
	config, err := util.LoadConfig(".")
	setMode(&config)
	if err != nil {
		log.Fatal().Err(err).Msg("config: cannot load config")
	}
	sqlDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("database: cannot connect to database")
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("database: ping database failed")
	}
	store := db.NewStore(sqlDB)
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	go runTaskProcessor(config, redisOpts, store)
	go runGRPCServer(config, store, taskDistributor)
	runGatewayServer(config, store, taskDistributor)
}

func runGatewayServer(config util.Config, store db.Store, distributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().Err(err).Msg("gateway: cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("gateway: cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("gateway: cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HTTPLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("gateway: cannot start HTTP gateway server")
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGRPCServer(config util.Config, store db.Store, distributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().Err(err).Msg("grpc: cannot create server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("net: unable to listen the port")
	}
	log.Printf("grpc: start grpc service at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("grpc: cannot start grpc server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("server: cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("server: cannot start server")
	}
}
