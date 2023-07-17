package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/Hypersus/simplebank/api"
	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/gapi"
	"github.com/Hypersus/simplebank/pb"
	"github.com/Hypersus/simplebank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgres://root:hypersus@localhost:5432/simple_bank?sslmode=disable"
// 	serverAddress = "localhost:8080"
// )

func main() {
	// gin.SetMode(gin.ReleaseMode)
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("config: cannot load config: ", err)
	}
	sqlDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("database: cannot connect to database: ", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("database: ping database failed: ", err)
	}
	store := db.NewStore(sqlDB)
	go runGRPCServer(config, store)
	runGatewayServer(config, store)
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("gateway: cannot create server: %s", err)
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
		log.Fatal("gateway: cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("gateway: cannot create listener: ", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("gateway: cannot start HTTP gateway server: ", err)
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("grpc: cannot create server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("net: unable to listen the port", err)
	}
	log.Printf("grpc: start grpc service at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("grpc: cannot start grpc server", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("server: cannot create server", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("server: cannot start server", err)
	}
}
