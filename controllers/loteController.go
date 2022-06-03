package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateLote(e echo.Context) (err error) {
	l := new(models.Lote)
	if err = e.Bind(l); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}\d]`, l.Lote)
	if f || len(l.Lote) > 4 || err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	lote := models.Lote{
		Lote: l.Lote,
	}
	if err = models.CreateLote(database.Ctx, &lote); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetLotes(e echo.Context) (err error) {
	locations := models.GetLotes(database.Ctx)
	return e.JSON(http.StatusOK, locations)
}
