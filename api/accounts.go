package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errMessage(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusConflict, errMessage(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errMessage(err))
		return
	}
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errMessage(err))
		return
	}
	authPayload, ok := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errMessage(errors.New("invalid auth payload")))
		return
	}
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errMessage(errors.New("not allowed")))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required, min=1"`
	PageSize int32 `form:"page_size" binding:"required, min=5, max=10"`
}

func (s *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errMessage(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
