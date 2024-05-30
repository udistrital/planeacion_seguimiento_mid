package helpers

import (
	"encoding/json"

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
