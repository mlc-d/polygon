package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateHistory(e echo.Context) (err error) {
	h := new(models.History)
	if err = e.Bind(h); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	history := models.History{
		ItemID:     h.ItemID,
		SkuID:      h.SkuID,
		LocationID: h.LocationID,
		StatusID:   h.StatusID,
		UserID:     h.UserID,
	}
	if err = models.CreateHistory(database.Ctx, &history); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetHistories(e echo.Context) (err error) {
	histories := models.GetHistories(database.Ctx)
	return e.JSON(http.StatusOK, histories)
}
