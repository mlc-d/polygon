package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateStatus(e echo.Context) (err error) {
	if !(utils.VerifyRole(e, 2)) {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	s := new(models.Status)
	if err = e.Bind(s); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}\d-]`, s.Status)
	if f || err != nil || len(s.Status) > 20 {
		err = errors.New(utils.Msg["invalidData"])
		return err
	}
	status := models.Status{
		Status: s.Status,
	}
	if err = models.CreateStatus(database.Ctx, &status); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetStatuses(e echo.Context) (err error) {
	statuses := models.GetStatuses(database.Ctx)
	return e.JSON(http.StatusOK, statuses)
}
