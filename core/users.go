package core

import (
	"log"
	"net/http"
	"rps/auth"
	"rps/models"
	"rps/models/db"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	Users []*models.User
)

// all but self

// @tags users
// @summary gives list of all registered users
// @description gives list of users excluding self
// @param X-Auth-Token header string true "token"
// @produce json
// @success 200 {array} models.User
// @failure 401 {string} string "auth required"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /users [get]
func GetUsers(c echo.Context) error {
	token, err := auth.Token(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	id, err := db.RedisClient.Get("SESSION:" + token).Result()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	users, err := db.GetUsers(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// just helper to get user id by x-auth from redis
func CurrentUserID(c echo.Context) (uint, error) {
	token, err := auth.Token(c)
	if err != nil {
		log.Printf("%+v", err)
		return 0, err
	}

	idStr, err := db.RedisClient.Get("SESSION:" + token).Result()
	if err != nil {
		log.Printf("%+v", err)
		return 0, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("%+v", err)
		return 0, err
	}
	return uint(id), nil
}
