package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateItem(c echo.Context) (err error) {
	i := new(models.Item)
	if err = c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	item := models.Item{
		Id:         i.Id,
		UIC:        i.UIC,
		LoteID:     i.LoteID,
		SkuID:      i.SkuID,
		LocationID: i.LocationID,
		StatusID:   i.StatusID,
		UserID:     i.UserID,
	}
	if err = models.CreateItem(database.Ctx, &item); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetItems(c echo.Context) (err error) {
	locations := models.GetItems(database.Ctx)
	return c.JSON(http.StatusOK, locations)
}

func UpdateItem(c echo.Context) (err error) {
	i := new(models.Item)
	if err = c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	// validate request
	if err = models.UpdateItem(database.Ctx, i); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusOK, utils.Response{
		"success": "actualizado",
	})
}

func AllocateItem(c echo.Context) (err error) {
	i := new(models.Item)
	if err = c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	// validate request
	if err = models.AllocateItem(database.Ctx, i); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusOK, utils.Response{
		"success": "reubicado",
	})
}
