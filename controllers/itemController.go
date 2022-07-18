package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateItem(c echo.Context) (err error) {
	items := new([]models.Item)
	if err = c.Bind(items); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}

	unchangedItems := utils.Response{}

	for _, item := range *items {
		if err = item.CreateItem(database.Ctx); err != nil {
			unchangedItems[err.Error()] = "cannot create item"
		}
	}

	if len(unchangedItems) > 0 {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": unchangedItems,
		})
	}

	return c.String(http.StatusCreated, fmt.Sprintf("ok - creado"))
}

func GetItems(c echo.Context) (err error) {
	locations := models.GetItems(database.Ctx)
	return c.JSON(http.StatusOK, locations)
}

func UpdateItem(c echo.Context) (err error) {
	i := new(models.Item)
	if err = c.Bind(i); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
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
	i := new([]models.Item)

	if err = c.Bind(i); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}

	unchangedItems := make([]string, 0)
	for _, v := range *i {
		if err = v.AllocateItem(database.Ctx); err != nil {
			unchangedItems = append(unchangedItems, err.Error()+" cannot reallocate item")
			// unchangedItems[err.Error()] = "cannot reallocate item"
		}
	}
	if len(unchangedItems) > 0 {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": unchangedItems,
		})
	}
	return c.String(http.StatusOK, "ok")
}

func GetItem(c echo.Context) (err error) {
	uic := c.QueryParam("uic")

	fmt.Println()
	fmt.Println(c.Request())
	fmt.Println()

	item, err := models.GetItem(database.Ctx, uic)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("error: %s", utils.Msg["dbError"]))
	}
	return c.JSON(http.StatusOK, item)
}
