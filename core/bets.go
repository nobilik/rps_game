package core

import (
	"net/http"
	"rps/models"
	"rps/models/db"

	"github.com/labstack/echo/v4"
)

type NewBetRequest struct {
	AcceptorID uint    `json:"acceptor_id"`
	Item       string  `json:"item"`
	Amount     float64 `json:"amount"`
}

type BetResponseRequest struct {
	BetID uint   `json:"acceptor_id"`
	Item  string `json:"item"`
}

// here we need amount, item & another player ID
// we check if user has enough money

// @tags bets
// @summary User place bet with amount, item name and another user ID
// @description New bet
// @param X-Auth-Token header string true "token"
// @param betRequest body NewBetRequest true "a new bet data"
// @accept json
// @success 200 {string} string "bet placed"
// @failure 400 {string} string "Some error"
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /bets [post]
func PlaceBet(c echo.Context) error {
	id, err := CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "auth required")
	}
	var data NewBetRequest
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Something went wrong with provided data")
	}

	item, ok := models.ItemsByName[data.Item]
	if !ok {
		return c.JSON(http.StatusBadRequest, "there is no such item")
	}
	err = db.Freeze(id, data.AcceptorID, data.Amount, *item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "bet placed")
}

// other users bets
// where acceptor is current user id

// @tags bets
// @summary gives list of active bets
// @description gives list of active bets that other users made with current user id as acceptor
// @param X-Auth-Token header string true "token"
// @produce json
// @success 200 {array} models.Bet
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /bets [get]
func GetBets(c echo.Context) error {
	id, err := CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "auth required")
	}

	bets, err := db.GetUserPendingBets(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bets)
}

// @tags bets
// @summary User place bet with amount, item name and another user ID
// @description New bet
// @param X-Auth-Token header string true "token"
// @param action path string true "decline or accept"
// @param bet body BetResponseRequest true "a bet response data"
// @accept json
// @success 200 {string} string "some result"
// @failure 400 {string} string "Some error"
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /bets/{action} [patch]
func RespondToBet(c echo.Context) error {
	userID, err := CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "auth required")
	}
	var data BetResponseRequest
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Something went wrong with provided data")
	}

	item, ok := models.ItemsByName[data.Item]
	if !ok {
		return c.JSON(http.StatusBadRequest, "there is no such item")
	}

	switch c.Param("action") {
	case "decline":
		err := db.Release(data.BetID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "declined")
	case "accept":
		bet, err := db.GetBetByID(data.BetID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "there is no such bet")
		}
		// check if acceptor has enough money to respond
		user, err := db.GetUserByID(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "there is no such user")
		}
		if user.Balance < bet.Amount {
			return c.JSON(http.StatusBadRequest, "insufficient funds")
		}
		// initialItem is the one that user selected on new bet
		initialItem, ok := models.ItemsByValue[bet.ItemValue]
		if !ok {
			return c.JSON(http.StatusBadRequest, "there is no such first item")
		}
		result := initialItem.Beats(item)

		// onLose && onWin from bet initiator side
		switch result {
		case models.Loss: // means that acceptor wins so we send notification (respond to acceptor)
			err := db.OnLose(bet, item)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, "you win")

		case models.Equal:
			err := db.OnEqual(bet, item)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, "nobody wins")
		case models.Winning:
			err := db.OnWin(bet, item)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, "you lose")
		default:
			return c.JSON(http.StatusInternalServerError, "something went wrong")
		}
	default:
		return c.JSON(http.StatusInternalServerError, "something went wrong")
	}
}
