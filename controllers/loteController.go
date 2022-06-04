package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateLote(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, 3)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	l := new(models.Lote)
	if err = c.Bind(l); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}\d]`, l.Lote)
	if f || len(l.Lote) > 4 || err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	lote := models.Lote{
		Lote: l.Lote,
	}
	if err = models.CreateLote(database.Ctx, &lote); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetLotes(c echo.Context) (err error) {
	locations := models.GetLotes(database.Ctx)
	return c.JSON(http.StatusOK, locations)
}
