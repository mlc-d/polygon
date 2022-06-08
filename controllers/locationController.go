package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateLocation(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, 4)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	l := new(models.Location)
	if err = c.Bind(l); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	location := models.Location{
		Location: l.Location,
		StatusID: l.StatusID,
	}
	if err = models.CreateLocation(database.Ctx, &location); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetLocations(c echo.Context) (err error) {
	locations := models.GetLocations(database.Ctx)
	return c.JSON(http.StatusOK, locations)
}
