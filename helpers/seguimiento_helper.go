package helpers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func CambiarEstadoPlan(plan map[string]interface{}, idEstado string) (map[string]interface{}, error) {
	idPlan := plan["_id"].(string)
	plan["estado_plan_id"] = idEstado
	var res map[string]interface{}
	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+idPlan, "PUT", &res, plan); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func ConvertirStringJson(diccionario map[string]interface{}) map[string]interface{} {
	dicStrings := map[string]interface{}{}
	for clave, valor := range diccionario {
		if clave == "informacion" || clave == "cualitativo" || clave == "cuantitativo" || clave == "estado" {
			datoJson := make(map[string]interface{})
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else if clave == "evidencia" {
			var datoJson []map[string]interface{}
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else {
			dicStrings[clave] = valor
		}
	}
	return dicStrings
}

// encodeBase62 genera el hash a partir de un planId en formato hexadecimal
func EncodeBase62(actividadID string) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	num := new(big.Int)
	num.SetString(actividadID, 16)

	var encoded string

	zero := big.NewInt(0)
	base := big.NewInt(62)

	for num.Cmp(zero) > 0 {
		var remainder big.Int
		num.DivMod(num, base, &remainder)
		encoded = string(charset[remainder.Int64()]) + encoded
	}

	return encoded
}

func DecodeBase62(base62Str string) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	decodedNum := new(big.Int)

	for _, char := range base62Str {
		index := strings.IndexByte(charset, byte(char))
		if index == -1 {
			panic("Invalid character found in base62 string")
		}
		decodedNum.Mul(decodedNum, big.NewInt(62))
		decodedNum.Add(decodedNum, big.NewInt(int64(index)))
	}

	return fmt.Sprintf("%x", decodedNum)
}
