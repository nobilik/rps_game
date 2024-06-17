package auth

import (
	"fmt"
	"net/http"
	"rps/helpers"
	"rps/models"
	"rps/models/db"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// @tags auth
// @summary Login with login-pass
// @description Gives value of X-Auth-Token for authentication
// @param LoginData body LoginData true "login-pass"
// @accept json
// @produce json
// @success 200 {object} LoginResponse "Successful"
// @failure 401 {string} string "bad login"
// @failure 500 {string} string "Something went wrong with provided data"
// @router       /login [post]
// Login with login-pass
func Login(c echo.Context) error {

	var data LoginData
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Something went wrong with provided data")
	}
	user, err := db.GetUserByLogin(data.Login)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "bad login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "bad login")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"token": createSession(user)})
}

func createSession(user *models.User) string {
	xAuthToken := helpers.GenerateWebSaveToken()
	// we can use here expiration
	db.RedisClient.Set("SESSION:"+xAuthToken, user.ID, 0)
	return xAuthToken
}

func Token(c echo.Context) (string, error) {
	payload := echo.Map{}
	err := (&echo.DefaultBinder{}).BindHeaders(c, &payload)
	if err != nil {
		return "", c.JSON(http.StatusBadRequest, "Something went wrong with provided data")
	}
	//
	// fmt.Printf("%+v", payload["X-Auth-Token"].([]string))
	tokenHeader, ok := payload["X-Auth-Token"].([]string)
	if !ok {
		return "", c.JSON(http.StatusForbidden, "no token provided")
	}
	fmt.Printf("%+v", tokenHeader[0])
	return tokenHeader[0], nil
}
