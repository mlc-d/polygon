package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateLocation(e echo.Context) (err error) {
	l := new(models.Location)
	if err = e.Bind(l); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	location := models.Location{
		Location: l.Location,
	}
	if err = models.CreateLocation(database.Ctx, &location); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetLocations(e echo.Context) (err error) {
	locations := models.GetLocations(database.Ctx)
	return e.JSON(http.StatusOK, locations)
}
