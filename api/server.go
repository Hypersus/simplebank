package api

import (
	"errors"
	"log"

	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/token"
	"github.com/Hypersus/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	ErrCustomValidator = errors.New("failed to register custom validator")
)

type Server struct {
	// data layer
	store db.Store
	// token maker for authentication
	tokenMaker token.TokenMaker
	// gin router
	router *gin.Engine
	// configuration
	config util.Config
}

func setRouter(server *Server) (err error) {
	if server.router != nil {
		return errors.New("router already set")
	}
	router := gin.Default()
	// Add the routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.POST("/transfers", server.createTransfer)
	server.router = router
	return nil
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	Maker := token.NewMaker(config.TokenType)
	if Maker == nil {
		return nil, token.ErrUnregisteredMaker
	}
	tokenMaker, err := Maker(config.TokenKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", currencyValidator); err != nil {
			return nil, ErrCustomValidator
		}
	} else {
		log.Fatal("gin: broken gin binding validator")
	}
	setRouter(server)
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errMessage(err error) gin.H {
	return gin.H{"error": err.Error()}
}
