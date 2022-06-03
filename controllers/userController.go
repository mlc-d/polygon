package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
	"strconv"
)

func CreateUser(e echo.Context) (err error) {
	u := new(models.User)
	if err = e.Bind(u); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}.]`, u.Name)
	if f || err != nil || len(u.Name) > 30 || len(u.Password) < 4 {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	user := models.User{
		Name:     u.Name,
		Password: u.Password,
		RoleID:   u.RoleID,
	}
	if err = models.CreateUser(database.Ctx, &user); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetUser(e echo.Context) (err error) {
	st := e.Param("id")
	var id int
	if id, err = strconv.Atoi(st); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	var user models.User
	if user, err = models.GetUser(database.Ctx, &models.User{Id: uint(id)}); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusOK, user)
}

func GetUsers(e echo.Context) (err error) {
	users := models.GetUsers(database.Ctx)
	return e.JSON(http.StatusOK, users)
}
