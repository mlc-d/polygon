package controllers

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	conf "gitlab.com/mlcprojects/wms/config"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
	"strings"
	"time"
)

var (
	config = conf.Cf
)

func Login(e echo.Context) (err error) {
	user := new(models.User)
	if err = e.Bind(&user); err != nil {
		return e.JSON(http.StatusInternalServerError, utils.Response{
			"error": "El servidor no reconoce la información enviada",
		})
	}
	if user.RoleID, err = models.ValidateUser(database.Ctx, user); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": "Credenciales inválidas",
		})
	}
	token, err := createAccessToken(user.Name, utils.StringValue(user.RoleID))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, utils.Response{
			"error": "Problema al crear un token de acceso",
		})
	}
	err = createRefreshToken(e, *user)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, utils.Response{
			"error": "Problema al crear token",
		})
	}
	return e.JSON(http.StatusOK, utils.Response{
		"accToken": token,
	})
}

func createAccessToken(subject, rolId string) (accessToken string, err error) {
	expireDate := time.Now().Add(time.Minute * 5)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.CustomJWTClaims{
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
		return nil, c.String(http.StatusInternalServerError, "")
	}
	return nil, nil
}

// Refresh Token

func createRefreshToken(e echo.Context, u models.User) (err error) {
	expireDate := time.Now().Add(time.Hour * 24 * 7)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &utils.CustomJWTClaims{
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
	e.SetCookie(cookie)
	return nil
}

func Refresh(c echo.Context) (err error) {
	requestCookie, err := c.Request().Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"content": err.Error(),
		})
	}
	requestToken := requestCookie.String()
	requestToken = strings.Split(requestToken, "refreshToken=")[1]
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
	if !refToken.Valid {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": utils.Msg["jwtError"],
		})
	}
	claims := refToken.Claims.(jwt.MapClaims)
	subject := utils.StringValue(claims["sub"])
	rolID := utils.StringValue(claims["rol"])
	newAccToken, err := createAccessToken(subject, rolID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{
			"error": utils.Msg["jwtError"],
		})
	}
	return c.JSON(http.StatusOK, utils.Response{
		"accToken": newAccToken,
	})
}
