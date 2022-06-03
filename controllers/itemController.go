package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateItem(e echo.Context) (err error) {
	i := new(models.Item)
	if err = e.Bind(i); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	item := models.Item{
		UIC:        i.UIC,
		SkuID:      i.SkuID,
		LocationID: i.LocationID,
		StatusID:   i.StatusID,
		UserID:     i.UserID,
	}
	if err = models.CreateItem(database.Ctx, &item); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetItems(e echo.Context) (err error) {
	locations := models.GetItems(database.Ctx)
	return e.JSON(http.StatusOK, locations)
}
