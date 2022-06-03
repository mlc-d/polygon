package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateRole(e echo.Context) (err error) {
	r := new(models.Role)
	if err = e.Bind(r); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	role := models.Role{
		Role:    r.Role,
		IsAdmin: r.IsAdmin,
	}
	if err = models.CreateRole(database.Ctx, &role); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetRoles(e echo.Context) (err error) {
	roles := models.GetRoles(database.Ctx)
	return e.JSON(http.StatusOK, roles)
}
