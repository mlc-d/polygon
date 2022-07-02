package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	conf "gitlab.com/mlcprojects/wms/config"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
)

var (
	config = conf.Cf
)

// Login validates credentials against database, and returns the proper tokens if everything is OK
func Login(c echo.Context) (err error) {
	// create a user object with the information sent by the client
	user := new(models.User)
	if err = c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": "El servidor no reconoce la información enviada",
		})
	}
	// verify the user object against the information in the database
	if user.RoleID, user.Id, err = models.ValidateUser(database.Ctx, user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": "Credenciales inválidas",
		})
	}
	// creates the access token
	token, err := createAccessToken(user.Name, utils.StringValue(user.RoleID), utils.StringValue(user.Id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": "Problema al crear un token de acceso",
		})
	}
	// creates the refresh token, passing down the context so the information can be added to the cookies
	err = createRefreshToken(c, *user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": "Problema al crear token de refresco",
		})
	}
	// everything is OK, returns access token
	return c.JSON(http.StatusOK, utils.Response{
		"accToken": token,
	})
}

// createAccessToken takes in the subject (username) role id and user id. Returns a jwt with those claims.
func createAccessToken(subject, rolId, userId string) (accessToken string, err error) {
	expireDate := time.Now().Add(time.Minute * 5)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.CustomJWTClaims{
		userId,
		rolId,
		jwt.StandardClaims{
			Subject:   subject,
			ExpiresAt: expireDate.Unix(),
		},
	})
	accToken, err := claims.SignedString([]byte(config.Jwt.AccSecKey))
	if err != nil {
		return "", err
	}
	return accToken, nil
}

// ValidateAccessToken ParseTokenFunc used by the Echo's JWT middleware. Checks if the token is valid.
func ValidateAccessToken(auth string, c echo.Context) (interface{}, error) {
	accToken, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, utils.ThrowErrorString("unexpected signing algorithm")
		}
		return []byte(config.Jwt.AccSecKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !accToken.Valid {
		return nil, errors.New("invalid token")
	}
	claims, _ := accToken.Claims.(jwt.MapClaims)
	// sets a new key-value pair containing the role id of the user making the request
	c.Set("roleFromReq", claims["rol"])
	return nil, nil
}

// createRefreshToken sets up a refresh token and add it to the response's cookies
func createRefreshToken(c echo.Context, u models.User) (err error) {
	expireDate := time.Now().Add(time.Hour * 24 * 7)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &utils.CustomJWTClaims{
		utils.StringValue(u.Id),
		utils.StringValue(u.RoleID),
		jwt.StandardClaims{
			Subject:   u.Name,
			ExpiresAt: expireDate.Unix(),
		},
	})
	refreshToken, err := claims.SignedString([]byte(config.Jwt.RefSecKey))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Expires = expireDate
	cookie.HttpOnly = true
	c.SetCookie(cookie)
	return nil
}

// Refresh validates if the refresh token is valid, and returns new access token
func Refresh(c echo.Context) (err error) {
	// parses cookie
	requestCookie, err := c.Request().Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"content": err.Error(),
		})
	}
	requestToken := requestCookie.String()
	requestToken = strings.Split(requestToken, "refreshToken=")[1]

	// parses token from cookie
	refToken, err := jwt.Parse(requestToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, utils.ThrowErrorString("unexpected signing algorithm")
		}
		return []byte(config.Jwt.RefSecKey), nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": utils.Msg["jwtError"],
		})
	}

	// checks if token is valid
	if !refToken.Valid {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": utils.Msg["jwtError"],
		})
	}

	// parses user's info and creates new access token
	claims := refToken.Claims.(jwt.MapClaims)
	subject := utils.StringValue(claims["sub"])
	rolID := utils.StringValue(claims["rol"])
	userID := utils.StringValue(claims["uid"])
	newAccToken, err := createAccessToken(subject, rolID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": utils.Msg["jwtError"],
		})
	}

	// everything is OK, returns new access token
	return c.JSON(http.StatusOK, utils.Response{
		"accToken": newAccToken,
	})
}

// Logout deletes refresh token cookie
func Logout(c echo.Context) (err error) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"ok": "false",
		})
	}
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, utils.Response{
		"ok": "logged out",
	})
}
