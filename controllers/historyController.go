package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateHistory(c echo.Context) (err error) {
	h := new(models.History)
	if err = c.Bind(h); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
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
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetHistories(c echo.Context) (err error) {
	histories := models.GetHistories(database.Ctx)
	return c.JSON(http.StatusOK, histories)
}

func GetItemHistory(c echo.Context) (err error) {
	uic := c.QueryParam("uic")
	var histories []models.PublicHistory
	if err, histories = models.GetItemHistory(database.Ctx, uic); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusOK, histories)
}
