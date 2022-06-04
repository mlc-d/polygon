package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateSku(e echo.Context) (err error) {
	if !(utils.VerifyRole(e, 4)) {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	s := new(models.Sku)
	if err = e.Bind(s); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["jsonError"],
		})
	}
	f, err := utils.ValidateInput(`[^\p{L}\d-]`, s.Sku)
	if f || err != nil || len(s.Sku) > 10 {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	sku := models.Sku{
		Sku:       s.Sku,
		ProductID: s.ProductID,
	}
	if err = models.CreateSku(database.Ctx, &sku); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetSkus(e echo.Context) (err error) {
	locations := models.GetLotes(database.Ctx)
	return e.JSON(http.StatusOK, locations)
}
