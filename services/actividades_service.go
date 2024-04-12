package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/planeacion"
	"github.com/udistrital/utils_oas/request"
)

func ConsultarActividadesGenerales(seguimiento_identificador string) ([]map[string]interface{}, error) {
	var respuestaSeguimiento map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	var respuesta map[string]interface{}
	var subgrupos []map[string]interface{}
	var seguimiento []map[string]interface{}
	var seguimientoDetalle []map[string]interface{}
	var datoPlan map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,_id:"+seguimiento_identificador, &respuestaSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimiento)

		if fmt.Sprintf("%v", seguimiento) != "[]" {
			planIdentificador := seguimiento[0]["plan_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planIdentificador, &respuesta); err == nil {
				request.LimpiezaRespuestaRefactor(respuesta, &subgrupos)

				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
						actividades := consultarActividades(subgrupos[i]["_id"].(string))

						if seguimiento[0]["dato"] == "{}" {
							for _, actividad := range actividades {
								actividad["estado"] = map[string]interface{}{"nombre": "Sin reporte"}
							}
						} else {
							dato_plan_str := seguimiento[0]["dato"].(string)
							json.Unmarshal([]byte(dato_plan_str), &datoPlan)

							for indiceActividad, elemento := range datoPlan {
								for _, actividad := range actividades {
									hayRegistro := false
									if reflect.TypeOf(actividad["index"]).String() == "string" {
										hayRegistro = indiceActividad == actividad["index"]
									} else {
										hayRegistro = indiceActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)
									}

									if hayRegistro {
										_, datosUnidos := elemento.(map[string]interface{})["estado"]
										if datosUnidos {
											actividad["estado"] = elemento.(map[string]interface{})["estado"]
											break
										} else {
											if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle?query=activo:true,_id:"+elemento.(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
												request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &seguimientoDetalle)
												dato := make(map[string]interface{})
												json.Unmarshal([]byte(seguimientoDetalle[0]["estado"].(string)), &dato)
												actividad["estado"] = dato
												break
											}
										}
									}
								}
							}
							for _, actividad := range actividades {
								if actividad["estado"] == nil {
									actividad["estado"] = map[string]interface{}{"nombre": "Sin reporte"}
								}
							}
						}
						return actividades, nil
					}
				}
			}
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	return nil, nil
}

func RevisarActividad(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (interface{}, error) {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	aux := make([]map[string]interface{}, 1)
	estado := map[string]interface{}{}
	comentario := false
	codigo_abreviacion := ""

	if errJSON := json.Unmarshal(requestBody, &body); errJSON != nil {
		logs.Error("Error -->", errJSON)
		return nil, errors.New(errJSON.Error())
	}
	if errPeticion := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); errPeticion != nil {
		logs.Error("Error -->", errPeticion)
		return nil, errors.New(errPeticion.Error())
	}

	request.LimpiezaRespuestaRefactor(respuesta, &aux)
	seguimiento = aux[0]
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)
	dato[indiceActividad] = body

	// Cualitativo
	if body["cualitativo"].(map[string]interface{})["observaciones"] != "" && body["cualitativo"].(map[string]interface{})["observaciones"] != "Sin observación" && body["cualitativo"].(map[string]interface{})["observaciones"] != nil {
		comentario = true
	}

	// Cuantitativo
	for _, indicador := range body["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{}) {
		if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" && indicador.(map[string]interface{})["observaciones"] != nil {
			comentario = true
			break
		}
	}

	// Evidencia
	for _, evidencia := range body["evidencia"].([]interface{}) {
		if evidencia.(map[string]interface{})["Observacion"] != "" && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" {
			comentario = true
			break
		}
	}

	if comentario {
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
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	identificador, segregado := body["id"].(string)
	if segregado && identificador != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+identificador, &respuestaDetalle); err == nil {
			request.LimpiezaRespuestaRefactor(respuestaDetalle, &detalle)
			detalle = planeacion.ConvertirStringJson(detalle)
			detalle["evidencia"] = body["evidencia"]
			detalle["cualitativo"] = body["cualitativo"]
			detalle["cuantitativo"] = body["cuantitativo"]
			detalle["estado"] = estado
			guardarDetalleSeguimiento(detalle, true)
		}
	} else {
		dato[indiceActividad].(map[string]interface{})["estado"] = estado
		b, _ := json.Marshal(dato)
		seguimiento["dato"] = string(b)
	}

	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	data := respuesta["Data"].(map[string]interface{})
	data["Observación"] = comentario
	return data, nil
}

func RetornarActividad(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (interface{}, error) {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(requestBody, &body); err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	aux := make([]map[string]interface{}, 1)
	request.LimpiezaRespuestaRefactor(respuesta, &aux)
	seguimiento = aux[0]
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)
	identificador, segregado := body["id"].(string)

	if segregado && identificador != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+identificador, &respuestaDetalle); err == nil {
			request.LimpiezaRespuestaRefactor(respuestaDetalle, &detalle)
			detalle = planeacion.ConvertirStringJson(detalle)
		} else {
			logs.Error("Error -->", err)
			return nil, errors.New(err.Error())
		}

		if detalle["estado"].(map[string]interface{})["id"] == "63793207242b813898e9856b" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:OAPC", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
			seguimiento["estado_seguimiento_id"] = estado["id"]
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
			detalle["estado"] = estado
			guardarDetalleSeguimiento(detalle, true)

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
			return respuesta["Data"].(map[string]interface{}), nil
		} else {
			return nil, nil
		}
	} else {
		dato[indiceActividad] = body

		if dato[indiceActividad].(map[string]interface{})["estado"].(map[string]interface{})["id"] == "63793207242b813898e9856b" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:OAPC", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}
			seguimiento["estado_seguimiento_id"] = estado["id"]
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
			dato[indiceActividad].(map[string]interface{})["estado"] = estado

			b, _ := json.Marshal(dato)
			seguimiento["dato"] = string(b)

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
			return respuesta["Data"].(map[string]interface{}), nil
		} else {
			return nil, nil
		}
	}
}
