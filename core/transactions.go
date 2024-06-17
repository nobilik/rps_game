package core

import (
	"net/http"
	"rps/models/db"

	"github.com/labstack/echo/v4"
)

type TransactionRequest struct {
	Amount float64 `json:"amount"`
}

// @tags transactions
// @summary gives list of own transactions
// @description gives list of own transactions
// @param X-Auth-Token header string true "token"
// @produce json
// @success 200 {array} models.Transaction
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /transactions [get]
func GetPersonalTransactions(c echo.Context) error {
	userID, err := CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "auth required")
	}
	transactions, err := db.GetTransactionsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}

// @tags transactions
// @summary Manages own balance
// @description top up or withdraw
// @param X-Auth-Token header string true "token"
// @param action path string true "action"
// @param transactionRequest body TransactionRequest true "object with amount"
// @accept json
// @success 200 {string} string "some result"
// @failure 400 {string} string "Some error"
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /transactions/{action} [post]
func ManageBalance(c echo.Context) error {
	userID, err := CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "auth required")
	}

	var data TransactionRequest
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Something went wrong with provided data")
	}
	switch c.Param("action") {
	case "topup":
		err := db.TopUp(userID, data.Amount)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "topped up")
	case "withdraw":
		err := db.Withdraw(userID, data.Amount)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "withdrawed")
	default:
		return c.JSON(http.StatusBadRequest, "bad action")

	}
}
