package services

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_seguimiento_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func GuardarDocumentos(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (map[string]interface{}, error) {
	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	codigo_abreviacion := ""
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(requestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Enlace"] != nil {
					evidencias = append(evidencias, evidencia.(map[string]interface{}))
					if evidencia.(map[string]interface{})["Observacion"] != nil && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" && evidencia.(map[string]interface{})["Observacion"] != "" {
						comentario = true
					}
				}
			}

			if body["documento"] != nil {
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
						"Observacion": "",
						"Activo":      true,
					})
				}
			}

			if body["unidad"].(bool) {
				codigo_abreviacion = "AER"
			} else if comentario {
				codigo_abreviacion = "CO"
			} else {
				codigo_abreviacion = "AR"
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}

			aux := make([]map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indiceActividad] == nil {
				detalle["evidencia"] = evidencias
				detalle["estado"] = estado
				delete(detalle, "_id")
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(detalle, false)}
			} else {
				identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]

				if segregado && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = helpers.ConvertirStringJson(detalle)
						detalle["evidencia"] = evidencias
						detalle["estado"] = estado
						helpers.GuardarDetalleSeguimiento(detalle, true)
					} else {
						logs.Error("Error --> ", err)
						return nil, errors.New(err.Error())
					}
				} else {
					dato[indiceActividad].(map[string]interface{})["evidencia"] = evidencias
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
				}
			}
			valor, _ := json.Marshal(dato)
			seguimiento["dato"] = string(valor)
			estadoSeguimiento, errEstadoSeg := helpers.ConsultarEstadoSeguimiento(seguimiento)
			if errEstadoSeg != nil {
				logs.Error("Error --> ", errEstadoSeg)
				return nil, errors.New(errEstadoSeg.Error())
			}
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error --> ", err)
				return nil, errors.New("error del servicio GuardarDocumentos:    Error guardado documentos del seguimiento \"seguimiento[\"_id\"].(string)\"" + err.Error())
			}

			return map[string]interface{}{"seguimiento": detalle["evidencia"], "estadoActividad": estado}, nil
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
}

func GuardarCualitativo(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (interface{}, error) {
	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cualitativo map[string]interface{}
	var informacion map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	var estado map[string]interface{}

	if err := json.Unmarshal(requestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]
			cualitativo = body["cualitativo"].(map[string]interface{})
			informacion = body["informacion"].(map[string]interface{})
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indiceActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				} else {
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
				}

				detalle = map[string]interface{}{"estado": estado, "cualitativo": cualitativo, "informacion": informacion}
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(detalle, false)}
			} else {
				identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]

				if segregado && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = helpers.ConvertirStringJson(detalle)
						estado = detalle["estado"].(map[string]interface{})
					} else {
						logs.Error("Error --> ", err)
						return nil, errors.New(err.Error())
					}
				} else {
					estado = dato[indiceActividad].(map[string]interface{})["estado"].(map[string]interface{})
				}

				if estado["nombre"] == "Con observaciones" && body["dependencia"].(bool) {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				} else if estado["nombre"] == "Actividad reportada" || estado["nombre"] == "Con observaciones" {
					var codigo_abreviacion string

					observacion = actividadConObservaciones(body)
					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AAV"
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					} else {
						logs.Error("Error --> ", err)
						return nil, errors.New(err.Error())
					}
				}

				if segregado && identificador != "" {
					detalle["cualitativo"] = cualitativo
					detalle["informacion"] = informacion
					detalle["estado"] = estado
					detalle["evidencia"] = body["evidencia"]
					helpers.GuardarDetalleSeguimiento(detalle, true)
				} else {
					dato[indiceActividad].(map[string]interface{})["informacion"] = informacion
					dato[indiceActividad].(map[string]interface{})["cualitativo"] = cualitativo
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
				}
			}
			b, _ := json.Marshal(dato)
			seguimiento["dato"] = string(b)
			estadoSeguimiento, errEstadoSeg := helpers.ConsultarEstadoSeguimiento(seguimiento)
			if errEstadoSeg != nil {
				logs.Error("Error --> ", errEstadoSeg)
				return nil, errors.New(errEstadoSeg.Error())
			}
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New("error del servicio GuardarCualitativo: Error actualizando componente cualitativo de seguimiento \"seguimiento[\"_id\"].(string)\"" + err.Error())
			}
			return respuesta["Data"], nil
		} else {
			logs.Error("Error -->", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func GuardarCuantitativo(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (interface{}, error) {
	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cuantitativo map[string]interface{}
	var informacion map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(requestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]
			cuantitativo = body["cuantitativo"].(map[string]interface{})
			informacion = body["informacion"].(map[string]interface{})
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indiceActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				} else {
					logs.Error("Error -->", err)
					return nil, errors.New(err.Error())
				}
				detalle = map[string]interface{}{"estado": estado, "cuantitativo": cuantitativo, "informacion": informacion}
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(detalle, false)}
			} else {
				identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]

				if segregado && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = helpers.ConvertirStringJson(detalle)
						estado = detalle["estado"].(map[string]interface{})
					}
				} else {
					estado = dato[indiceActividad].(map[string]interface{})["estado"].(map[string]interface{})
				}

				if estado["nombre"] == "Con observaciones" && body["dependencia"].(bool) {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				} else if estado["nombre"] == "Actividad reportada" || estado["nombre"] == "Con observaciones" {
					var codigo_abreviacion string

					observacion = actividadConObservaciones(body)
					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AAV"
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}

				if segregado && identificador != "" {
					detalle["estado"] = estado
					detalle["cuantitativo"] = cuantitativo
					detalle["informacion"] = informacion
					helpers.GuardarDetalleSeguimiento(detalle, true)
				} else {
					dato[indiceActividad].(map[string]interface{})["informacion"] = informacion
					dato[indiceActividad].(map[string]interface{})["cuantitativo"] = cuantitativo
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
				}
			}

			b, _ := json.Marshal(dato)
			seguimiento["dato"] = string(b)
			estadoSeguimiento, errEstadoSeg := helpers.ConsultarEstadoSeguimiento(seguimiento)
			if errEstadoSeg != nil {
				logs.Error("Error --> ", errEstadoSeg)
				return nil, errors.New(errEstadoSeg.Error())
			}
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New("error del servicio GuardarCuantitativo:    Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"" + err.Error())
			}

			return respuesta["Data"], nil
		} else {
			logs.Error("Error -->", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func actividadConObservaciones(seguimiento map[string]interface{}) bool {
	var cuantitativo map[string]interface{}
	var cualitativo map[string]interface{}

	if seguimiento["cuantitativo"] != nil {
		cuantitativo = seguimiento["cuantitativo"].(map[string]interface{})
		for _, indicador := range cuantitativo["indicadores"].([]interface{}) {
			if (indicador.(map[string]interface{})["observaciones_dependencia"] != "" && indicador.(map[string]interface{})["observaciones_dependencia"] != "Sin observación" && indicador.(map[string]interface{})["observaciones_dependencia"] != nil) ||
				(indicador.(map[string]interface{})["observaciones_planeacion"] != "" && indicador.(map[string]interface{})["observaciones_planeacion"] != "Sin observación" && indicador.(map[string]interface{})["observaciones_planeacion"] != nil) {
				return true
			}
		}
	}

	if seguimiento["cualitativo"] != nil {
		cualitativo = seguimiento["cualitativo"].(map[string]interface{})
		if (cualitativo["observaciones_dependencia"] != "" && cualitativo["observaciones_dependencia"] != "Sin observación" && cualitativo["observaciones_dependencia"] != nil) ||
			(cualitativo["observaciones_planeacion"] != "" && cualitativo["observaciones_planeacion"] != "Sin observación" && cualitativo["observaciones_planeacion"] != nil) {
			return true
		}
	}

	if seguimiento["evidencia"] != nil {
		for _, evidencia := range seguimiento["evidencia"].([]map[string]interface{}) {
			if evidencia["Observacion"] != "" && evidencia["Observacion"] != "Sin observación" && evidencia["Observacion"] != nil {
				return true
			}
		}
	}

	return false
}
