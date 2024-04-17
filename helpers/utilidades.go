package helpers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/planeacion"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func GuardarDocumento(documentos []interface{}) []interface{} {
	var respuestaDocumentos []interface{}

	for _, documento := range documentos {

		if documento.(map[string]interface{})["file"] != nil {
			documento := map[string]interface{}{
				"IdTipoDocumento": documento.(map[string]interface{})["IdTipoDocumento"],
				"nombre":          documento.(map[string]interface{})["nombre"],
				"metadatos":       documento.(map[string]interface{})["metadatos"],
				"descripcion":     documento.(map[string]interface{})["descripcion"],
				"file":            documento.(map[string]interface{})["file"],
			}

			var documentoAuxiliar []map[string]interface{}
			documentoAuxiliar = append(documentoAuxiliar, documento)
			documentoSubido, errorDocumento := RegistrarDocumento(documentoAuxiliar)

			if errorDocumento == nil {
				documentoTemporal := map[string]interface{}{
					"Nombre":        documentoSubido.(map[string]interface{})["Nombre"].(string),
					"Enlace":        documentoSubido.(map[string]interface{})["Enlace"],
					"Id":            documentoSubido.(map[string]interface{})["Id"],
					"TipoDocumento": documentoSubido.(map[string]interface{})["TipoDocumento"],
					"Activo":        documentoSubido.(map[string]interface{})["Activo"],
				}

				respuestaDocumentos = append(respuestaDocumentos, documentoTemporal)
			}
		}
	}
	return respuestaDocumentos
}

func RegistrarDocumento(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}
	errRegDoc := request.SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"/document/uploadAnyFormat", "POST", &resultadoRegistro, documento)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}

}

func GuardarDetalleSeguimiento(detalle map[string]interface{}, actualizar bool) string {
	var respuesta map[string]interface{}
	var identificador string
	var tipo_peticion string
	url := "http://" + beego.AppConfig.String("PlanesService") + "/seguimiento-detalle"

	detalle = planeacion.ConvertirJsonString(detalle)
	data := map[string]interface{}{}

	data["estado"] = detalle["estado"].(string)
	if fmt.Sprintf("%v", detalle["activo"]) == "<nil>" {
		data["activo"] = "true"
	} else {
		data["activo"] = detalle["activo"].(string)
	}
	if fmt.Sprintf("%v", detalle["fecha_creacion"]) == "<nil>" {
		data["fecha_creacion"] = time_bogota.TiempoBogotaFormato()
		data["fecha_modificacion"] = time_bogota.TiempoBogotaFormato()
	} else {
		data["fecha_creacion"] = detalle["fecha_creacion"].(string)
		data["fecha_modificacion"] = detalle["fecha_modificacion"].(string)
	}

	if valor, existe := detalle["informacion"]; existe {
		detalle["informacion"] = valor
	} else {
		detalle["informacion"] = "{}"
	}
	if valor, existe := detalle["cualitativo"]; existe {
		detalle["cualitativo"] = valor
	} else {
		detalle["cualitativo"] = "{}"
	}
	if valor, existe := detalle["cuantitativo"]; existe {
		detalle["cuantitativo"] = valor
	} else {
		detalle["cuantitativo"] = "{}"
	}
	if valor, existe := detalle["evidencia"]; existe {
		detalle["evidencia"] = valor
	} else {
		detalle["evidencia"] = "[]"
	}

	if actualizar {
		tipo_peticion = "PUT"
		url += "/" + detalle["_id"].(string)
	} else {
		tipo_peticion = "POST"
	}

	if err := request.SendJson(url, tipo_peticion, &respuesta, data); err == nil && respuesta["Status"].(string) != "400" {
		aux := make(map[string]interface{})
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		identificador = aux["_id"].(string)
	}
	return identificador
}

func ConsultarEstadoSeguimiento(seguimiento map[string]interface{}) (string, error) {
	var respuestaEstado map[string]interface{}
	enReporte := true
	var estado map[string]interface{}
	dato := make(map[string]interface{})
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &respuestaEstado); err == nil {
		estado = map[string]interface{}{
			"nombre": respuestaEstado["Data"].(map[string]interface{})["nombre"],
			"id":     respuestaEstado["Data"].(map[string]interface{})["_id"],
		}

		for _, actividad := range dato {
			_, datosUnidos := actividad.(map[string]interface{})["estado"]

			if datosUnidos {
				if actividad.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad en reporte" {
					enReporte = false
				}
			} else {
				var respuestaSeguimientoDetalle map[string]interface{}
				var seguimientoDetalle map[string]interface{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+actividad.(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
					request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &seguimientoDetalle)
					dato := make(map[string]interface{})
					json.Unmarshal([]byte(seguimientoDetalle["estado"].(string)), &dato)
					if dato["nombre"] != "Actividad en reporte" {
						enReporte = false
					}
				} else {
					logs.Error("Error --> ", err)
					return "", errors.New(err.Error())
				}
			}
		}

		if enReporte {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:ER", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}
		}
	} else {
		logs.Error("Error --> ", err)
		return "", errors.New(err.Error())
	}
	return estado["id"].(string), nil
}
