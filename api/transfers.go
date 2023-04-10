package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateTransferParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1,nefield=FromAccountID"`
	Amount        int64  `json:"amount" binding:"required,min=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errMessage(err))
		return
	}
	if !s.isValidTransfer(ctx, req.FromAccountID, req.Currency) || !s.isValidTransfer(ctx, req.ToAccountID, req.Currency) {
		return
	}
	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := s.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) isValidTransfer(ctx *gin.Context, id int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errMessage(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errMessage(err))
		return false
	}
	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch: account currency is %v, transfer currency is %s", account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errMessage(err))
		return false
	}
	return true
}
