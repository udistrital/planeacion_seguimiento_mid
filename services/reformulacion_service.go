package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_seguimiento_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func ObtenerPeriodoUltimoReporte(planId string) (map[string]interface{}, string, error) {
	var resPlan map[string]interface{}
	var plan map[string]interface{}
	indiceUltimoTrimestreAvalado := -1

	// Get plan
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+planId, &resPlan); err != nil {
		return nil, "", errors.New("error del servicio ObtenerPeriodoUltimoReporte: Error al consultar el plan")
	}
	request.LimpiezaRespuestaRefactor(resPlan, &plan)

	// Obtener trimestres de la vigencia
	trimestres, err := ConsultarTrimestres(plan["vigencia"].(string))
	if err != nil || len(trimestres) == 0 {
		return nil, "", errors.New("error del servicio ObtenerPeriodoUltimoReporte: No se encontraron trimestres")
	}

	for i, trimestre := range trimestres {
		seguimiento, err := ConsultarEstadoTrimestre(planId, trimestre["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string))
		if err == nil {
			if fmt.Sprintf("%v", seguimiento.(map[string]interface{})["dato"]) != "{}" && seguimiento.(map[string]interface{})["estado_seguimiento_id"].(map[string]interface{})["codigo_abreviacion"].(string) == "AV" {
				if indiceUltimoTrimestreAvalado == i-1 {
					indiceUltimoTrimestreAvalado = i
				} else {
					return nil, "", errors.New("los seguimientos de los trimestres no han sido diligenciados en el orden correcto")
				}
			}
		} else {
			return nil, "", errors.New("error del servicio ObtenerPeriodoUltimoReporte: Error al consultar el seguimiento del trimestre")
		}
	}

	if indiceUltimoTrimestreAvalado == 3 {
		return nil, "", errors.New("no hay trimestres a los cuales se les pueda aplicar la reformulación")
	}
	// Guarda el trimestre desde el cual hará el siguiente seguimiento
	indiceUltimoTrimestreAvalado = indiceUltimoTrimestreAvalado + 1
	return plan, trimestres[indiceUltimoTrimestreAvalado]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string), nil
}

func SolicitarReformulacion(requestBody []byte) (map[string]interface{}, error) {
	var body map[string]interface{}
	var estadoFormulado []map[string]interface{}
	var respuesta map[string]interface{}
	var evidencias []map[string]interface{}
	var resultadoRegistro map[string]interface{}
	var registro map[string]interface{}

	json.Unmarshal(requestBody, &body)

	// Obtiene el estado 'Reformulación de Plan de Acción - Formulado - SisgPlan'
	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/?query=CodigoAbreviacion:RPA-F-SP&fields=Id", &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &estadoFormulado)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	// Obtiene el último periodo que haya sido reportado
	_, periodo, err := ObtenerPeriodoUltimoReporte(body["plan_id"].(string))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Guarda el/los documento(s)
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

	// Arma el cuerpo de la petición
	cuerpoPeticionCRUD := map[string]interface{}{
		"observaciones": body["observaciones"].(string),
		"plan_id":       body["plan_id"].(string),
		"periodo":       periodo,
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

func ValidacionReformulacion(planIdentificador string) (map[string]interface{}, error) {
	plan, periodo, err := ObtenerPeriodoUltimoReporte(planIdentificador)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	periodoSalida := map[string]interface{}{
		"plan":    plan,
		"periodo": periodo,
	}
	return periodoSalida, nil
}

func AprobarReformulacion(reformulacionId string) (map[string]interface{}, error) {
	var reformulacion map[string]interface{}
	var estadoAprobado []map[string]interface{}
	var respuesta map[string]interface{}
	var resultadoModificacionReformulacion map[string]interface{}
	var resultadoClonacionPlan map[string]interface{}
	var cuerpoPeticionClonacionFormato map[string]interface{}
	var modificacionReformulacion map[string]interface{}
	var resultadoModificacionPlanPadre map[string]interface{}

	// Obtiene la reformulacion
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/reformulacion/"+reformulacionId, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &reformulacion)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	// Obtiene el último periodo del trimestre que tuvo valores ingresados en seguimiento
	planPadre, _, err := ObtenerPeriodoUltimoReporte(reformulacion["plan_id"].(string))
	if err != nil {
		return nil, errors.New("error del servicio SolicitarReformulacion: Error al obtener el último seguimiento")
	}
	cuerpoPeticionClonacionFormato = make(map[string]interface{})
	cuerpoPeticionClonacionFormato["nombre"] = "" + planPadre["nombre"].(string)
	cuerpoPeticionClonacionFormato["descripcion"] = planPadre["descripcion"].(string)
	cuerpoPeticionClonacionFormato["tipo_plan_id"] = planPadre["tipo_plan_id"].(string)
	cuerpoPeticionClonacionFormato["aplicativo_id"] = planPadre["aplicativo_id"].(string)
	cuerpoPeticionClonacionFormato["activo"] = planPadre["activo"]
	cuerpoPeticionClonacionFormato["reformulacion"] = true
	if err := request.SendJson("http://"+beego.AppConfig.String("FormulacionService")+"/formulacion/formato/"+reformulacion["plan_id"].(string)+"/clonar", "POST", &resultadoClonacionPlan, cuerpoPeticionClonacionFormato); err == nil {
		request.LimpiezaRespuestaRefactor(resultadoClonacionPlan, &cuerpoPeticionClonacionFormato)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	// Modifica el estado activo del plan padre
	cuerpoPlanPadre := make(map[string]interface{})
	for k, v := range planPadre {
		cuerpoPlanPadre[k] = v
	}
	cuerpoPlanPadre["activo"] = false
	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+reformulacion["plan_id"].(string), "PUT", &resultadoModificacionPlanPadre, cuerpoPlanPadre); err == nil {
		request.LimpiezaRespuestaRefactor(resultadoModificacionPlanPadre, &cuerpoPlanPadre)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	// Modifica la reformulacion

	// Obtiene el estado de Reformulación de Plan de Acción Aprobado
	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/?query=CodigoAbreviacion:RPA-A-SP&fields=Id", &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &estadoAprobado)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	cuerpoReformulacion := make(map[string]interface{})
	for k, v := range reformulacion {
		cuerpoReformulacion[k] = v
	}
	cuerpoReformulacion["estado_id"] = estadoAprobado[0]["Id"].(int)

	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/reformulacion/", "PUT", &resultadoModificacionReformulacion, cuerpoReformulacion); err == nil {
		request.LimpiezaRespuestaRefactor(resultadoModificacionReformulacion, &modificacionReformulacion)
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	return modificacionReformulacion, nil
}
