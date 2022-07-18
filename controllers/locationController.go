package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
)

func CreateLocation(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, SUPERVISOR_ROLE_ID)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	l := new(models.Location)
	if err = c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
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

func GetLocation(c echo.Context) (err error) {
	loc := c.QueryParam("l")
	location, err := models.GetLocation(database.Ctx, loc)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.Response{
			"error": "ubicación no encontrada",
		})
	}
	return c.JSON(http.StatusOK, location)
}

func GetLocationByStatus(c echo.Context) (err error) {
	sid := c.QueryParam("status")
	statusId, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	locations, err := models.GetLocationByStatus(database.Ctx, uint(statusId))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, locations)
}

func EditLocation(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, LEADER_ROLE_ID)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	l := new(models.Location)
	if err = c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	location := models.Location{
		Id:       l.Id,
		Location: l.Location,
		StatusID: l.StatusID,
	}
	if err = models.EditLocation(database.Ctx, &location); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}
