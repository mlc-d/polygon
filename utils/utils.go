package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"regexp"
	"strconv"
)

type Response map[string]any

type CustomJWTClaims struct {
	Rid string `json:"rol"`
	jwt.StandardClaims
}

var (
	//Msg is a dictionary for possible json responses, in Spanish
	//keys:
	//jsonError
	//invalidData
	//dbError
	Msg = map[string]string{
		"jsonError":    "El servidor no reconoce la información enviada",
		"invalidData":  "La información enviada tiene inconsistencias",
		"dbError":      "La petición a la base de datos no tuvo éxito",
		"jwtError":     "Error al validar o generar un JWT",
		"unauthorized": "Su usuario no está autorizado para realizar la operación",
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

func ByteValue(i interface{}) []byte {
	return []byte(fmt.Sprintf("%v", i))
}

func ThrowErrorString(i interface{}) error {
	return fmt.Errorf("error: %v", i)
}

func VerifyRole(e echo.Context, downTo int) bool {
	rol := e.Get("rolFromReq")
	if val, err := strconv.Atoi(StringValue(rol)); err != nil || val >= downTo {
		return false
	}
	return true
}
