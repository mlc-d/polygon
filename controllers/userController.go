package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
	"strconv"
)

func CreateUser(c echo.Context) (err error) {
	u := new(models.User)

	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["jsonError"],
		})
	}
	if !(utils.VerifyRole(c, 4)) {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["unauthorized"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}.]`, u.Name)
	if f || err != nil || len(u.Name) > 30 || len(u.Password) < 4 {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["invalidData"],
		})
	}
	user := models.User{
		Name:     u.Name,
		Password: u.Password,
		RoleID:   u.RoleID,
	}
	if err = models.CreateUser(database.Ctx, &user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Res{
		"success": "creado",
	})
}

func GetUser(c echo.Context) (err error) {
	st := c.Param("id")
	var id int
	if id, err = strconv.Atoi(st); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["invalidData"],
		})
	}
	var user models.User
	if user, err = models.GetUser(database.Ctx, &models.User{Id: uint(id)}); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusOK, user)
}

func GetUsers(c echo.Context) (err error) {
	users := models.GetUsers(database.Ctx)
	return c.JSON(http.StatusOK, users)
}
