package controllers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
)

func CreateProduct(e echo.Context) (err error) {
	if !(utils.VerifyRole(e, 4)) {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}
	p := new(models.Product)
	if err = e.Bind(p); err != nil {
		return e.JSON(http.StatusBadRequest, "invalid entry, please check your request")
	}
	wrongName, err := utils.ValidateInput(`[^\p{L}\d.¡!#]`, p.Name)
	wrongRef, err1 := utils.ValidateInput(`[^\p{L}\d.:-_]`, p.Ref)
	wrongDescription, err2 := utils.ValidateInput(`[^\p{L}\d.,:;¡!#]`, p.Description)
	if wrongName || wrongRef || wrongDescription || err != nil || err1 != nil || err2 != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	product := models.Product{
		Name:        p.Name,
		Ref:         p.Ref,
		Description: p.Description,
	}
	if err = models.CreateProduct(database.Ctx, &product); err != nil {
		return e.JSON(http.StatusBadRequest, utils.Response{
			"error": "dbError",
		})
	}
	return e.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetProducts(e echo.Context) (err error) {
	products := models.GetProducts(database.Ctx)
	return e.JSON(http.StatusOK, products)
}
