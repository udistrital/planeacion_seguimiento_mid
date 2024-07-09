package services

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_seguimiento_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func SolicitarReformulacion(requestBody []byte) (map[string]interface{}, error) {
	var body map[string]interface{}
	var estadoFormulado []map[string]interface{}
	var respuesta map[string]interface{}
	var evidencias []map[string]interface{}
	var resultadoRegistro map[string]interface{}
	var registro map[string]interface{}

	json.Unmarshal(requestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/?query=CodigoAbreviacion:RPA-F-SP&fields=Id", &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &estadoFormulado)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	respuestaDocs := helpers.GuardarDocumento(body["documento"].([]interface{}))
	for _, doc := range respuestaDocs {

		evidencias = append(evidencias, map[string]interface{}{
			"Id":     doc.(map[string]interface{})["Id"],
			"Enlace": doc.(map[string]interface{})["Enlace"],
			"nombre": doc.(map[string]interface{})["Nombre"],
			"TipoDocumento": map[string]interface{}{
				"id":                doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"],
				"codigoAbreviacion": doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["CodigoAbreviacion"],
			},
			"Observacion": body["observaciones"].(string),
			"Activo":      true,
		})
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{"documentos": evidencias})
	if err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	archivos := string(jsonBytes)

	cuerpoPeticionCRUD := map[string]interface{}{
		"observaciones": body["observaciones"].(string),
		"plan_id":       body["plan_id"].(string),
		"estado_id":     estadoFormulado[0]["Id"].(float64),
		"archivos":      archivos,
	}

	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/reformulacion/", "POST", &resultadoRegistro, cuerpoPeticionCRUD); err == nil {
		request.LimpiezaRespuestaRefactor(resultadoRegistro, &registro)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	return registro, nil
}
