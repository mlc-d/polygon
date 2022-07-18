package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateBatch(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, 3)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	l := new(models.Batch)
	if err = c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	f, err := utils.ValidateInput(`[^\p{L}\d]`, l.Batch)
	if f || len(l.Batch) > 4 || err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	batch := models.Batch{
		Batch: l.Batch,
	}
	if err = models.CreateBatch(database.Ctx, &batch); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetBatches(c echo.Context) (err error) {
	locations := models.GetBatches(database.Ctx)
	return c.JSON(http.StatusOK, locations)
}
