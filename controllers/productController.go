package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateProduct(c echo.Context) (err error) {
	if !(utils.VerifyRole(c, 4)) {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["unauthorized"],
		})
	}
	p := new(models.Product)
	if err = c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["jsonError"],
		})
	}
	/*wrongName, err := utils.ValidateInput(`[^\p{L}\d.¡!# ]`, p.Name)
	wrongRef, err1 := utils.ValidateInput(`[^\p{L}\d.:-_]`, p.Ref)
	wrongDescription, err2 := utils.ValidateInput(`[^\p{L}\d.,:;¡!# ]`, p.Description)
	if wrongName || wrongRef || wrongDescription || err != nil || err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": utils.Msg["invalidData"],
		})
	}*/
	product := models.Product{
		Name:        p.Name,
		Ref:         p.Ref,
		Description: p.Description,
	}
	if err = models.CreateProduct(database.Ctx, &product); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Res{
			"error": "dbError",
		})
	}
	return c.JSON(http.StatusCreated, utils.Res{
		"success": "creado",
	})
}

func GetProducts(c echo.Context) (err error) {
	products := models.GetProducts(database.Ctx)
	return c.JSON(http.StatusOK, products)
}
