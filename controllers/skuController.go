package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateSku(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, 4)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	s := new(models.Sku)
	if err = c.Bind(s); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	/*f, err := utils.ValidateInput(`[^\p{L}\d-]`, s.Sku)
	if f || err != nil || len(s.Sku) > 10 {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}*/
	sku := models.Sku{
		Sku:       s.Sku,
		ProductID: s.ProductID,
	}
	if err = models.CreateSku(database.Ctx, &sku); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetSkus(c echo.Context) (err error) {
	skus := models.GetSkus(database.Ctx)
	return c.JSON(http.StatusOK, skus)
}
