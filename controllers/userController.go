package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/utils"
	"net/http"
	"strconv"
	"strings"
)

func CreateUser(c echo.Context) (err error) {
	u := new(models.User)

	if err = c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	if !(utils.VerifyRole(c, MANAGER_ROLE_ID)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}

	u.Name = sanitizeUserName(u.Name)
	notOk, err := utils.ValidateInput(`[^\p{L}.]`, u.Name)
	if notOk || err != nil || len(u.Name) > 30 || len(u.Password) < 4 {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	user := models.User{
		Name:     u.Name,
		Password: u.Password,
		RoleID:   u.RoleID,
	}
	if err = models.CreateUser(database.Ctx, &user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		"success": "creado",
	})
}

func GetUser(c echo.Context) (err error) {
	st := c.Param("id")
	var id int
	if id, err = strconv.Atoi(st); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	var user models.User
	if user, err = models.GetUser(database.Ctx, &models.User{Id: uint(id)}); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["dbError"],
		})
	}
	return c.JSON(http.StatusOK, user)
}

func GetUsers(c echo.Context) (err error) {
	users := models.GetUsers(database.Ctx)
	return c.JSON(http.StatusOK, users)
}

func EditUser(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}

	if !utils.VerifyRole(c, int(u.RoleID)) {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["unauthorized"],
		})
	}

	// sanitizes data
	u.Name = sanitizeUserName(u.Name)
	ok, err := utils.ValidateInput(`[^\p{L}.]`, u.Name)
	if !ok || err != nil || len(u.Name) > 30 {
		return c.JSON(http.StatusBadRequest, utils.Response{
			"error": utils.Msg["invalidData"],
		})
	}
	user := models.User{
		Id:       u.Id,
		Name:     u.Name,
		Password: u.Password,
	}
	if err = models.UpdateUser(database.Ctx, &user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["dbError"]))
	}
	return c.String(http.StatusCreated, "ok - creado")
}

func DeleteUser(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["jsonError"]))
	}
	err = models.DeleteUser(database.Ctx, u)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", utils.Msg["dbError"]))
	}

	return c.String(http.StatusCreated, "ok - eliminado")
}

func sanitizeUserName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
