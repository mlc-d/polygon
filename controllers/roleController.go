package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateRole(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, ADMIN_ROLE_ID)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	r := new(models.Role)
	if err = c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	role := models.Role{
		Role:    r.Role,
		IsAdmin: r.IsAdmin,
	}
	if err = models.CreateRole(database.Ctx, &role); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetRoles(c echo.Context) (err error) {
	roles := models.GetRoles(database.Ctx)
	return c.JSON(http.StatusOK, roles)
}
