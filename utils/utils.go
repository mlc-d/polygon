package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"regexp"
	"strconv"
)

type Response map[string]interface{}

type CustomJWTClaims struct {
	Uid string `json:"uid"`
	Rid string `json:"rol"`
	jwt.StandardClaims
}

var (
	// Msg is a dictionary for possible json responses, in Spanish
	// keys:
	// jsonError
	// invalidData
	// dbError
	// jwtError
	// unauthorized
	Msg = map[string]string{
		"jsonError":          "El servidor no reconoce la información enviada",
		"invalidCredentials": "Credenciales Inválidas",
		"invalidData":        "La información enviada tiene inconsistencias",
		"dbError":            "La petición a la base de datos no tuvo éxito",
		"jwtError":           "Error al validar o generar un JWT",
		"unauthorized":       "Su usuario no está autorizado para realizar la operación",
	}
)

// ValidateInput Simple wrapper for regexp.Match(), just converting input
// to the proper type
func ValidateInput(pattern string, input string) (bool, error) {
	return regexp.Match(pattern, []byte(input))
}

// StringValue Simple wrapper for fmt.Sprintf(),
func StringValue(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

// ByteValue converts input to string, and then to byte array
// func ByteValue(i interface{}) []byte {
//	return []byte(fmt.Sprintf("%v", i))
// }

// VerifyRole each request sends the user's role id. Controllers take it and
// compare it to the list of roles authorized to perform the action.
func VerifyRole(c echo.Context, downTo int) bool {
	rol := c.Get("roleFromReq")
	if val, err := strconv.Atoi(StringValue(rol)); err != nil || val >= downTo {
		return false
	}
	return true
}
