package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func ObtenerPeriodoUltimoReporte(planId string) (arrPeriodo map[string]interface{}, errRes error) {
	var resPlan map[string]interface{}
	var plan map[string]interface{}
	var resSeguimientos map[string]interface{}
	var seguimientos []map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var resPeriodos map[string]interface{}
	var periodos []map[string]interface{}
	var seguimientosLlenos []map[string]interface{}
	var seguimientosVacios []map[string]interface{}

	// Get plan
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+planId, &resPlan); err != nil {
		errRes = errors.New("error del servicio AvalarPlan: Error al consultar plan")
		return nil, errRes
	}
	request.LimpiezaRespuestaRefactor(resPlan, &plan)

	// Obtener trimestres de la vigencia
	trimestres, err := ConsultarTrimestres(plan["vigencia"].(string))
	if err != nil || len(trimestres) == 0 {
		errRes = errors.New("error del servicio AvalarPlan: No se encontraron trimestres")
		return nil, errRes
	}
	formatdata.JsonPrint(trimestres)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId, &resSeguimientos); err == nil {
		request.LimpiezaRespuestaRefactor(resSeguimientos, &seguimientos)
		for _, seguimiento := range seguimientos {
			if fmt.Sprintf("%v", seguimiento["dato"]) != "{}" {
				seguimientosLlenos = append(seguimientosLlenos, seguimiento)
			} else {
				seguimientosVacios = append(seguimientosVacios, seguimiento)
			}
		}
	}
	// fmt.Println("--------------------SLlenos------------")
	// formatdata.JsonPrint(seguimientosLlenos)
	// fmt.Println("--------------------SVacios------------")
	// formatdata.JsonPrint(seguimientosVacios)

	for _, seguimientoActual := range seguimientosLlenos {
		// fmt.Println("http://" + beego.AppConfig.String("PlanesService") + "/periodo-seguimiento/" + seguimientoActual["periodo_seguimiento_id"].(string))
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimientoActual["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
			request.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)
			fmt.Println("--------------------periodoSeguimiento------------")
			fmt.Println(periodoSeguimiento["periodo_id"])
			// formatdata.JsonPrint(periodoSeguimiento)
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodos); err == nil {
				fmt.Println("--------------------periodo------------")
				request.LimpiezaRespuestaRefactor(resPeriodos, &periodos)
				formatdata.JsonPrint(periodos[0])

			}

		}

	}

	return trimestres[0], nil
}

func SolicitarReformulacion(requestBody []byte) (map[string]interface{}, error) {
	var body map[string]interface{}
	// var estadoFormulado []map[string]interface{}
	// var respuesta map[string]interface{}
	// var evidencias []map[string]interface{}
	// var resultadoRegistro map[string]interface{}
	var registro map[string]interface{}

	json.Unmarshal(requestBody, &body)

	// if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/?query=CodigoAbreviacion:RPA-F-SP&fields=Id", &respuesta); err == nil {
	// 	request.LimpiezaRespuestaRefactor(respuesta, &estadoFormulado)
	// } else {
	// 	logs.Error("Error -->", err)
	// 	return nil, errors.New(err.Error())
	// }

	// respuestaDocs := helpers.GuardarDocumento(body["documento"].([]interface{}))
	// for _, doc := range respuestaDocs {

	// 	evidencias = append(evidencias, map[string]interface{}{
	// 		"Id":     doc.(map[string]interface{})["Id"],
	// 		"Enlace": doc.(map[string]interface{})["Enlace"],
	// 		"nombre": doc.(map[string]interface{})["Nombre"],
	// 		"TipoDocumento": map[string]interface{}{
	// 			"id":                doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"],
	// 			"codigoAbreviacion": doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["CodigoAbreviacion"],
	// 		},
	// 		"Observacion": body["observaciones"].(string),
	// 		"Activo":      true,
	// 	})
	// }

	// jsonBytes, err := json.Marshal(map[string]interface{}{"documentos": evidencias})
	// if err != nil {
	// 	logs.Error("Error -->", err)
	// 	return nil, errors.New(err.Error())
	// }
	// archivos := string(jsonBytes)

	// cuerpoPeticionCRUD := map[string]interface{}{
	// 	"observaciones": body["observaciones"].(string),
	// 	"plan_id":       body["plan_id"].(string),
	// 	"estado_id":     estadoFormulado[0]["Id"].(float64),
	// 	"archivos":      archivos,
	// }

	// if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/reformulacion/", "POST", &resultadoRegistro, cuerpoPeticionCRUD); err == nil {
	// 	request.LimpiezaRespuestaRefactor(resultadoRegistro, &registro)
	// } else {
	// 	logs.Error("Error -->", err)
	// 	return nil, errors.New(err.Error())
	// }
	ObtenerPeriodoUltimoReporte(body["plan_id"].(string))
	registro = make(map[string]interface{})
	registro["plan_id"] = body["plan_id"].(string)
	return registro, nil
}
