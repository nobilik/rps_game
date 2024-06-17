package auth

import (
	"log"
	"net/http"
	"rps/models"
	"rps/models/db"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// @tags auth
// @summary Register with login-pass
// @description Adds user to db
// @param RegisterData body LoginData true "login-pass"
// @accept json
// @success 201 {string} string "created"
// @failure 400 {string} string "Some error"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /register [post]
// Register with login-pass
func Register(c echo.Context) error {
	var data LoginData
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Something went wrong with provided data")
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Unable to encode pass: %s", data.Password)
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	user := &models.User{
		Login:    data.Login,
		Password: string(encrypted),
	}
	err = db.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "created")
}
