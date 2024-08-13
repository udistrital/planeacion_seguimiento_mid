package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	comunhelper "github.com/udistrital/planeacion_mid/helpers/comunHelper"
	"github.com/udistrital/planeacion_seguimiento_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"golang.org/x/sync/errgroup"
)

func GuardarSeguimiento(requestBody []byte, planIdentificador string, indiceActividad string, trimestre string) (interface{}, error) {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	var estado map[string]interface{}

	if err := json.Unmarshal(requestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Enlace"] != nil {
					evidencias = append(evidencias, evidencia.(map[string]interface{}))
				}
			}
			body["evidencia"] = evidencias
			aux := make([]map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			} else {
				logs.Error("Error --> ", err)
				return nil, errors.New(err.Error())
			}

			if dato[indiceActividad] == nil {
				body["estado"] = estado
				delete(body, "_id")
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(body, false)}
				valor, _ := json.Marshal(dato)
				str := string(valor)
				seguimiento["dato"] = str
			} else {
				identificador, actualizar := dato[indiceActividad].(map[string]interface{})["id"].(string)

				if actualizar && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = helpers.ConvertirStringJson(detalle)
						detalle["estado"] = estado
						helpers.GuardarDetalleSeguimiento(detalle, true)
					}
				} else {
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
					valor, _ := json.Marshal(dato)
					str := string(valor)
					seguimiento["dato"] = str
				}
			}
			estadoSeguimiento, errEstadoSeg := helpers.ConsultarEstadoSeguimiento(seguimiento)
			if errEstadoSeg != nil {
				logs.Error("Error --> ", errEstadoSeg)
				return nil, errors.New(errEstadoSeg.Error())
			}
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error --> ", err)
				return nil, errors.New("error del servicio GuardarSeguimiento:    Error actualizando seguimiento \"seguimiento[\"_id\"].(string)\"" + err.Error())
			}

			return respuesta["Data"], nil
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
}

func ConsultarSeguimiento(planIdentificador string, indiceActividad string, trimestreIdentificador string) (map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var respuestaPeriodoSeguimiento map[string]interface{}
	var respuestaPeriodo map[string]interface{}
	var respuestaEstado map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimientoActividad map[string]interface{}
	var periodo []map[string]interface{}
	var trimestre string
	dato := make(map[string]interface{})

	id_actividad_decoded := planIdentificador + "" + indiceActividad
	id_actividad := helpers.EncodeBase62(id_actividad_decoded)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestreIdentificador, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &respuestaPeriodoSeguimiento); err == nil {
			request.LimpiezaRespuestaRefactor(respuestaPeriodoSeguimiento, &periodoSeguimiento)

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &respuestaPeriodo); err == nil {
				request.LimpiezaRespuestaRefactor(respuestaPeriodo, &periodo)
				trimestre = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)
			} else {
				logs.Error("Error --> ", err)
				return nil, errors.New(err.Error())
			}
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)
		actividad, _ := json.Marshal(consultarActividad(seguimiento, indiceActividad, trimestre))
		json.Unmarshal([]byte(string(actividad)), &seguimientoActividad)
		seguimientoActividad["_id"] = seguimiento["_id"].(string)
		seguimientoActividad["id_actividad"] = id_actividad

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &respuestaEstado); err == nil {
			seguimientoActividad["estadoSeguimiento"] = respuestaEstado["Data"].(map[string]interface{})["nombre"].(string)
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}
		return seguimientoActividad, nil
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
}

func RevisarSeguimiento(seguimientoIdentificador string) (interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,_id:"+seguimientoIdentificador, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)
		avalado, observacion, mensaje, errAvalable := seguimientoAvalable(seguimiento)

		if (avalado || observacion) && errAvalable == nil {
			var codigo_abreviacion string

			if observacion {
				codigo_abreviacion = "CO"
			} else {
				codigo_abreviacion = "AV"
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &respuestaEstado); err == nil {
				estado := map[string]interface{}{
					"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
				seguimiento["estado_seguimiento_id"] = estado["id"]
			}

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}

			data := respuesta["Data"].(map[string]interface{})
			return data, nil
		} else {
			if mensaje != nil {
				responseJSON, _ := json.Marshal(mensaje)
				return nil, fmt.Errorf("%v", string(responseJSON))
			} else {
				logs.Error("Error -->", errAvalable)
				return nil, errors.New(errAvalable.Error())
			}
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func MigrarInformacion(planIdentificador string, trimestre string) (interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	respuestaMigrado := []map[string]interface{}{}
	respuestaNoMigrado := []map[string]interface{}{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)

		for indiceActividad, actividad := range dato {
			identificador, segregado := actividad.(map[string]interface{})["id"].(string)

			if !segregado || identificador == "" {
				delete(dato[indiceActividad].(map[string]interface{}), "_id")
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(dato[indiceActividad].(map[string]interface{}), false)}
				respuestaMigrado = append(respuestaMigrado, map[string]interface{}{"id": indiceActividad})
			} else {
				respuestaNoMigrado = append(respuestaNoMigrado, map[string]interface{}{"id": indiceActividad})
			}
			b, _ := json.Marshal(dato)
			seguimiento["dato"] = string(b)

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New("error del servicio MigrarInformacion:    Error actualizando el seguimiento nuevo" + err.Error())
			}
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	return map[string]interface{}{"Actividades migradas:": respuestaMigrado, "Actividades no migradas: ": respuestaNoMigrado}, nil
}

func VerificarSeguimiento(idSeguimiento string) (interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+idSeguimiento, &respuesta); err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
	aux := make(map[string]interface{}, 1)
	request.LimpiezaRespuestaRefactor(respuesta, &aux)
	seguimiento = aux

	if reportable, mensaje := seguimientoReportable(seguimiento); reportable == false {
		logs.Error("Error -->", reportable)
		return nil, fmt.Errorf("%v", mensaje)
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:RV", &resEstado); err == nil {
		seguimiento["estado_seguimiento_id"] = resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
	}

	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	return seguimiento, nil
}

func consultarActividad(seguimiento map[string]interface{}, indice string, trimestre string) map[string]interface{} {
	var data map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var informacion map[string]interface{}
	var cuantitativo map[string]interface{}
	cualitativo := map[string]interface{}{}
	evidencia := []interface{}{}
	evidenciaSeg := []map[string]interface{}{}
	estado := map[string]interface{}{}
	detalle := map[string]interface{}{}
	identificador := ""
	dato := make(map[string]interface{})
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[indice] != nil {
		identificadores, segregado := dato[indice].(map[string]interface{})["id"]

		if segregado && identificadores != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+identificadores.(string), &respuestaDetalle); err == nil {
				if respuestaDetalle["Data"] != "null" {
					request.LimpiezaRespuestaRefactor(respuestaDetalle, &detalle)
					detalle = helpers.ConvertirStringJson(detalle)
					identificador = detalle["_id"].(string)

					if len(detalle["informacion"].(map[string]interface{})) == 0 {
						informacion, _ = consultarInformacionPlan(seguimiento, indice)
					} else {
						informacion = detalle["informacion"].(map[string]interface{})
					}

					if len(detalle["cuantitativo"].(map[string]interface{})) == 0 {
						cuantitativo, _ = consultarCuantitativoPlan(seguimiento, indice, trimestre)
					} else {
						cuantitativo = detalle["cuantitativo"].(map[string]interface{})
					}

					if len(detalle["cualitativo"].(map[string]interface{})) == 0 {
						cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
					} else {
						cualitativo = detalle["cualitativo"].(map[string]interface{})
					}

					if len(detalle["evidencia"].([]map[string]interface{})) != 0 {
						evidenciaSeg = detalle["evidencia"].([]map[string]interface{})
					}

					if len(detalle["estado"].(map[string]interface{})) == 0 {
						if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &respuestaEstado); err == nil {
							estado = map[string]interface{}{
								"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
								"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
							}
						}
					} else {
						estado = detalle["estado"].(map[string]interface{})
					}
				}
			}
		} else {
			if dato[indice].(map[string]interface{})["informacion"] == nil {
				informacion, _ = consultarInformacionPlan(seguimiento, indice)
			} else {
				informacion = dato[indice].(map[string]interface{})["informacion"].(map[string]interface{})
			}

			if dato[indice].(map[string]interface{})["cuantitativo"] == nil {
				cuantitativo, _ = consultarCuantitativoPlan(seguimiento, indice, trimestre)
			} else {
				cuantitativo = dato[indice].(map[string]interface{})["cuantitativo"].(map[string]interface{})
			}

			if dato[indice].(map[string]interface{})["evidencia"] != nil {
				evidencia = dato[indice].(map[string]interface{})["evidencia"].([]interface{})
			}

			if dato[indice].(map[string]interface{})["estado"] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else {
				estado = dato[indice].(map[string]interface{})["estado"].(map[string]interface{})
			}

			if dato[indice].(map[string]interface{})["cualitativo"] == nil {
				cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
			} else {
				cualitativo = dato[indice].(map[string]interface{})["cualitativo"].(map[string]interface{})
			}
		}
	} else {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &respuestaEstado); err == nil {
			estado = map[string]interface{}{
				"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
				"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
			}
		}
		informacion, _ = consultarInformacionPlan(seguimiento, indice)
		cuantitativo, _ = consultarCuantitativoPlan(seguimiento, indice, trimestre)
		cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
	}

	data = map[string]interface{}{
		"id":           identificador,
		"informacion":  informacion,
		"cualitativo":  cualitativo,
		"cuantitativo": cuantitativo,
		"estado":       estado,
		"evidencia":    evidencia,
	}

	if identificador != "" {
		data["evidencia"] = evidenciaSeg
	}

	return data
}

func consultarInformacionPlan(seguimiento map[string]interface{}, indice string) (map[string]interface{}, error) {
	var respuestaPlan map[string]interface{}
	var respuestaPeriodoSeguimiento map[string]interface{}
	var respuestaPeriodo map[string]interface{}
	var respuestaInformacion map[string]interface{}
	var respuestaDependencia []map[string]interface{}
	var hijos []map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var periodo []map[string]interface{}

	informacion := map[string]interface{}{
		"ponderacion": "",
		"periodo":     "",
		"tarea":       "",
		"producto":    "",
		"nombre":      "",
		"descripcion": "",
		"index":       indice,
		"unidad":      "",
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
		informacion["nombre"] = respuestaPlan["Data"].(map[string]interface{})["nombre"]
		informacion["unidad"] = respuestaPlan["Data"].(map[string]interface{})["dependencia_id"]
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &respuestaPeriodoSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaPeriodoSeguimiento, &periodoSeguimiento)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &respuestaPeriodo); err == nil {
			request.LimpiezaRespuestaRefactor(respuestaPeriodo, &periodo)
			informacion["trimestre"] = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &respuestaInformacion); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaInformacion, &hijos)

		for _, hijo := range hijos {
			nombreHijo := strings.ToLower(hijo["nombre"].(string))

			if hijo["activo"] == true {
				var respuesta map[string]interface{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijo["_id"].(string), &respuesta); err == nil {
					datoPlan := make(map[string]interface{})
					dato := make(map[string]interface{})
					nombreDetalle := strings.ToLower(respuesta["Data"].([]interface{})[0].(map[string]interface{})["nombre"].(string))

					if strings.Contains(nombreDetalle, "indicadores") || strings.Contains(nombreDetalle, "indicador") {
						continue
					}

					json.Unmarshal([]byte(respuesta["Data"].([]interface{})[0].(map[string]interface{})["dato"].(string)), &dato)
					if dato["required"] == false || dato["required"] == "false" {
						continue
					}

					json.Unmarshal([]byte(respuesta["Data"].([]interface{})[0].(map[string]interface{})["dato_plan"].(string)), &datoPlan)
					if datoPlan[indice] == nil {
						continue
					}

					switch {
					case strings.Contains(nombreHijo, "ponderación"):
						informacion["ponderacion"] = datoPlan[indice].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "periodo") || strings.Contains(nombreHijo, "período"):
						informacion["periodo"] = datoPlan[indice].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "tareas") || strings.Contains(nombreHijo, "actividades específicas"):
						informacion["tarea"] = datoPlan[indice].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "producto"):
						informacion["producto"] = datoPlan[indice].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "actividad general"):
						informacion["descripcion"] = datoPlan[indice].(map[string]interface{})["dato"]
						continue
					}
				} else {
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
				}
			}
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+informacion["unidad"].(string), &respuestaDependencia); err == nil {
		informacion["unidad"] = respuestaDependencia[0]["DependenciaId"].(map[string]interface{})["Nombre"]
	} else {
		informacion["unidad"] = nil
	}
	return informacion, nil
}

func consultarCuantitativoPlan(seguimiento map[string]interface{}, indice string, trimestre string) (map[string]interface{}, error) {
	var respuestaInformacion map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var hijos []interface{}
	var subgrupos []map[string]interface{}
	var indicadores []map[string]interface{}
	respuestas := make([]map[string]interface{}, 0)
	response := map[string]interface{}{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &respuestaInformacion); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaInformacion, &subgrupos)

		for _, subgrupo := range subgrupos {
			if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") && subgrupo["activo"] == true {
				hijos = subgrupo["hijos"].([]interface{})
				hijos = append(hijos, subgrupo["_id"])

				for _, hijo := range hijos {
					var respuesta map[string]interface{}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijo.(string), &respuesta); err == nil {
						hijosIndicadores := respuesta["Data"].(map[string]interface{})["hijos"].([]interface{})
						var dato_plan map[string]interface{}
						informacion := map[string]interface{}{
							"detalleReporte": "",
						}
						respuesta := map[string]interface{}{
							"indicador":            0,
							"indicadorAcumulado":   0,
							"avanceAcumulado":      0,
							"brechaExistente":      0,
							"acumuladoNumerador":   0,
							"acumuladoDenominador": 0,
							"meta":                 0,
						}

						for _, hijoI := range hijosIndicadores {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijoI.(string), &respuestaDetalle); err == nil {
								var subgrupo_detalle []map[string]interface{}
								request.LimpiezaRespuestaRefactor(respuestaDetalle, &subgrupo_detalle)

								if len(subgrupo_detalle) > 0 {
									if subgrupo_detalle[0]["dato_plan"] != nil {
										dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
										json.Unmarshal([]byte(dato_plan_str), &dato_plan)
										nombreDetalle := strings.ToLower(subgrupo_detalle[0]["nombre"].(string))

										if dato_plan[indice] == nil || dato_plan[indice].(map[string]interface{})["dato"] == "" {
											break
										}

										switch {
										case strings.Contains(nombreDetalle, "nombre"):
											informacion["nombre"] = dato_plan[indice].(map[string]interface{})["dato"]
											respuesta["nombre"] = dato_plan[indice].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "meta"):
											informacion["meta"] = dato_plan[indice].(map[string]interface{})["dato"]
											if reflect.TypeOf(dato_plan[indice].(map[string]interface{})["dato"]).String() == "string" {
												respuesta["meta"], _ = strconv.ParseFloat(dato_plan[indice].(map[string]interface{})["dato"].(string), 64)
											} else {
												respuesta["meta"] = dato_plan[indice].(map[string]interface{})["dato"].(float64)
											}
											continue
										case strings.Contains(nombreDetalle, "fórmula"):
											informacion["formula"] = dato_plan[indice].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "criterio"):
											informacion["denominador"] = dato_plan[indice].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "tendencia"):
											informacion["tendencia"] = strings.Trim(dato_plan[indice].(map[string]interface{})["dato"].(string), " ")
											continue
										case strings.Contains(nombreDetalle, "unidad de medida"):
											informacion["unidad"] = strings.Trim(dato_plan[indice].(map[string]interface{})["dato"].(string), " ")
											respuesta["unidad"] = strings.Trim(dato_plan[indice].(map[string]interface{})["dato"].(string), " ")
											continue
										}
									}
								}
							} else {
								logs.Error("Error --> ", err)
								return nil, errors.New(err.Error())
							}
						}

						if informacion["reporteDenominador"] == 1.0 {
							informacion["reporteDenominador"] = nil
						}

						if informacion["nombre"] != nil && informacion["nombre"] != "" {
							indicadores = append(indicadores, informacion)
							respuestas = append(respuestas, respuesta)
						}

						respuestasAnteriores, errorRespuestaAnterior := consultarRespuestaAnterior(seguimiento, len(indicadores)-1, respuestas, indice, trimestre)
						if errorRespuestaAnterior != nil {
							logs.Error("Error --> ", errorRespuestaAnterior)
							return nil, errors.New(errorRespuestaAnterior.Error())
						} else {
							respuestas = respuestasAnteriores
						}
					} else {
						logs.Error("Error --> ", err)
						return nil, errors.New(err.Error())
					}
				}
				break
			}
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}

	response["indicadores"] = indicadores
	response["resultados"] = respuestas
	return response, nil
}

func seguimientoAvalable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false
	estado := map[string]interface{}{}

	planIdentificador := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planIdentificador, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				actividades, errActividades := ConsultarActividades(subgrupos[i]["_id"].(string))
				if errActividades == nil {
					for indiceActividad, elemento := range datoPlan {
						identificador, segregado := elemento.(map[string]interface{})["id"]

						for _, actividad := range actividades {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								if indiceActividad == actividad["index"] {
									if segregado && identificador != "" {
										if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+elemento.(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
											request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
											detalle = helpers.ConvertirStringJson(detalle)
											estado = detalle["estado"].(map[string]interface{})
											if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
												dato[indiceActividad] = actividad["dato"]
											}
										}
									} else {
										if elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
											dato[indiceActividad] = actividad["dato"]
										}
									}
								}
							} else if indiceActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
								if segregado && identificador != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+elemento.(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
										request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
										detalle = helpers.ConvertirStringJson(detalle)
										estado = detalle["estado"].(map[string]interface{})
										if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
											dato[indiceActividad] = actividad["dato"]
										}
									} else {
										logs.Error("Error -->", err)
										return avaladas, observaciones, nil, errors.New(err.Error())
									}
								} else {
									if elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
										dato[indiceActividad] = actividad["dato"]
									}
								}
							}
						}
						if segregado && identificador != "" {
							if estado["nombre"] == "Con observaciones" {
								observaciones = true
							}

							if estado["nombre"] == "Actividad avalada" {
								avaladas = true
							}
						} else {
							if elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Con observaciones" {
								observaciones = true
							}

							if elemento.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Actividad avalada" {
								avaladas = true
							}
						}
					}
				}
			}
		}
	} else {
		logs.Error("Error -->", err)
		return avaladas, observaciones, nil, errors.New(err.Error())
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return avaladas, observaciones, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}, nil
	}

	return avaladas, observaciones, nil, nil
}

func consultarRespuestaAnterior(dataSeg map[string]interface{}, indice int, respuestas []map[string]interface{}, indiceActividad string, trimestre string) ([]map[string]interface{}, error) {
	plan_identificador := dataSeg["plan_id"].(string)
	var respuestaSeguimiento map[string]interface{}
	var respuestaPeriodoSeguimiento map[string]interface{}
	var respuestaPeriodo map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodo []map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_identificador, &respuestaSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)

		acumuladoNumerador := 0.0
		acumuladoDenominador := 0.0
		indicadorAcumulado := 0.0
		avanceAcumulado := 0.0
		brechaExistente := 0.0
		divisionCero := false

		for _, seguimiento := range seguimientos {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &respuestaPeriodoSeguimiento); err == nil {
				request.LimpiezaRespuestaRefactor(respuestaPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &respuestaPeriodo); err == nil {
					request.LimpiezaRespuestaRefactor(respuestaPeriodo, &periodo)
					tri, _ := strconv.Atoi(string(trimestre[1]))
					segTrimestre, _ := strconv.Atoi(string(periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)[1]))

					if (tri - 1) == segTrimestre {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)

							if dato[indiceActividad] == nil {
								respuestas[indice]["indicadorAcumulado"] = indicadorAcumulado
								respuestas[indice]["avanceAcumulado"] = avanceAcumulado
								respuestas[indice]["brechaExistente"] = brechaExistente
								respuestas[indice]["divisionCero"] = divisionCero
								continue
							}

							identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]
							if segregado && identificador != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
									request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
									detalle = helpers.ConvertirStringJson(detalle)

									if fmt.Sprintf("%v", detalle["cuantitativo"]) == "map[]" {
										respuestas[indice]["indicadorAcumulado"] = indicadorAcumulado
										respuestas[indice]["avanceAcumulado"] = avanceAcumulado
										respuestas[indice]["brechaExistente"] = brechaExistente
										respuestas[indice]["divisionCero"] = divisionCero
										continue
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["indicadorAcumulado"] != nil {
										indicadorAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["indicadorAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["avanceAcumulado"] != nil {
										avanceAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["avanceAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["brechaExistente"] != nil {
										brechaExistente += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["brechaExistente"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["divisionCero"] != nil {
										divisionCero = detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["divisionCero"].(bool)
									} else {
										divisionCero = false
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoDenominador"] != nil {
										acumuladoDenominador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoDenominador"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoNumerador"] != nil {
										acumuladoNumerador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoNumerador"].(float64)
									}
								} else {
									logs.Error("Error --> ", err)
									return nil, errors.New(err.Error())
								}
							} else {
								seguimientoActividad := dato[indiceActividad].(map[string]interface{})
								if seguimientoActividad["cuantitativo"] == nil {
									respuestas[indice]["indicadorAcumulado"] = indicadorAcumulado
									respuestas[indice]["avanceAcumulado"] = avanceAcumulado
									respuestas[indice]["brechaExistente"] = brechaExistente
									respuestas[indice]["divisionCero"] = divisionCero
									continue
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["indicadorAcumulado"] != nil {
									indicadorAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["indicadorAcumulado"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["avanceAcumulado"] != nil {
									avanceAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["avanceAcumulado"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["brechaExistente"] != nil {
									brechaExistente += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["brechaExistente"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["divisionCero"] != nil {
									divisionCero = seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["divisionCero"].(bool)
								} else {
									divisionCero = false
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoDenominador"] != nil {
									acumuladoDenominador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoDenominador"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoNumerador"] != nil {
									acumuladoNumerador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[indice].(map[string]interface{})["acumuladoNumerador"].(float64)
								}
							}
						}

						respuestas[indice]["indicadorAcumulado"] = indicadorAcumulado
						respuestas[indice]["avanceAcumulado"] = avanceAcumulado
						respuestas[indice]["brechaExistente"] = brechaExistente
						respuestas[indice]["acumuladoNumerador"] = acumuladoNumerador
						respuestas[indice]["acumuladoDenominador"] = acumuladoDenominador
						respuestas[indice]["divisionCero"] = divisionCero
						break
					}
				} else {
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
				}
			} else {
				logs.Error("Error --> ", err)
				return nil, errors.New(err.Error())
			}
		}
	}
	return respuestas, nil
}

func EstadoTrimestres(planId string) (interface{}, error) {
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var auxPlanes []map[string]interface{}
	var mutex sync.Mutex
	wge := new(errgroup.Group)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId, &resSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(resSeguimiento, &planes)

		for _, plan := range planes {
			plan := plan
			wge.Go(func() error {
				var periodo []map[string]interface{}
				periodoSeguimientoId := plan["periodo_seguimiento_id"].(string)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=_id:"+periodoSeguimientoId, &resPeriodoSeguimiento); err == nil {
					var periodoSeguimiento []map[string]interface{}
					request.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)
					if fmt.Sprintf("%v", periodoSeguimiento[0]) != "map[]" {

						if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento[0]["periodo_id"].(string), &resPeriodo); err == nil {
							request.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
							auxPeriodo := periodo[0]["ParametroId"].(map[string]interface{})
							periodoSeguimiento[0]["periodo_nombre"] = auxPeriodo["CodigoAbreviacion"].(string)
							plan["periodo_seguimiento_id"] = periodoSeguimiento[0]

							if fmt.Sprintf("%v", periodo[0]) != "map[]" {
								var resEstado map[string]interface{}

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+plan["estado_seguimiento_id"].(string), &resEstado); err == nil {
									plan["estado_seguimiento_id"] = resEstado["Data"]

									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["plan_id"].(string), &resEstado); err == nil {
										plan["plan_id"] = resEstado["Data"]

										mutex.Lock()
										auxPlanes = append(auxPlanes, plan)
										mutex.Unlock()
									}
								}
							}
						} else {
							logs.Error("Error -->", err)
							return errors.New(err.Error())
						}
					}
				} else {
					logs.Error("Error -->", err)
					return errors.New(err.Error())
				}
				return nil
			})
		}

		if err := wge.Wait(); err != nil {
			return nil, errors.New(err.Error())
		}

		if auxPlanes != nil {
			// Ordenar auxPlanes por periodo_nombre antes de devolver
			sort.SliceStable(auxPlanes, func(i, j int) bool {
				return auxPlanes[i]["periodo_seguimiento_id"].(map[string]interface{})["periodo_nombre"].(string) <
					auxPlanes[j]["periodo_seguimiento_id"].(map[string]interface{})["periodo_nombre"].(string)
			})
			return auxPlanes, nil
		} else {
			return []int{}, nil
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func ConsultarEstadoTrimestre(planIdentificador string, trimestre string) (map[string]interface{}, error) {
	var respuestaSeguimiento map[string]interface{}
	var respuestaPeriodoSeguimiento map[string]interface{}
	var respuestaPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var periodoSeguimiento []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador, &respuestaSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &planes)

		for _, plan := range planes {
			var periodo []map[string]interface{}
			periodoSeguimientoIdentificador := plan["periodo_seguimiento_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=_id:"+periodoSeguimientoIdentificador, &respuestaPeriodoSeguimiento); err == nil {
				request.LimpiezaRespuestaRefactor(respuestaPeriodoSeguimiento, &periodoSeguimiento)
				if fmt.Sprintf("%v", periodoSeguimiento[0]) != "map[]" {

					if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento[0]["periodo_id"].(string)+",ParametroId__CodigoAbreviacion:"+trimestre, &respuestaPeriodo); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaPeriodo, &periodo)
						plan["periodo_seguimiento_id"] = periodoSeguimiento[0]

						if fmt.Sprintf("%v", periodo[0]) != "map[]" {
							var respuestaEstado map[string]interface{}

							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+plan["estado_seguimiento_id"].(string), &respuestaEstado); err == nil {
								plan["estado_seguimiento_id"] = respuestaEstado["Data"]

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["plan_id"].(string), &respuestaEstado); err == nil {
									plan["plan_id"] = respuestaEstado["Data"]

									return plan, nil
								}
							}
						}
					}
				}
			}
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
	return nil, nil
}

func AvalarPlan(plan_id string) (arrReportes []map[string]interface{}, errRes error) {
	var resPlan map[string]interface{}
	var plan map[string]interface{}
	id_estado_avalado := "6153355601c7a2365b2fb2a1"
	id_estado_preaval := "614d3b4401c7a222052fac05"

	// Get plan
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_id, &resPlan); err != nil {
		errRes = errors.New("error del servicio AvalarPlan: Error al consultar plan")
		return nil, errRes
	}
	if resPlan["Data"] == nil {
		errRes = errors.New("error del servicio AvalarPlan: Plan no encontrado")
		return nil, errRes
	}
	request.LimpiezaRespuestaRefactor(resPlan, &plan)

	// Cambiar estado plan a avalado
	respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_avalado)
	if err != nil || respuesta["Success"] == false {
		errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
		return nil, errRes
	}

	// Obtener trimestres de la vigencia
	trimestres, err := ConsultarTrimestres(plan["vigencia"].(string))
	if err != nil || len(trimestres) == 0 {
		errRes = errors.New("error del servicio AvalarPlan: No se encontraron trimestres")
		return nil, errRes
	}

	// Creacion de reportes de seguimiento
	tipo := "61f236f525e40c582a0840d0"
	var resPadres map[string]interface{}
	var resDependencia []map[string]interface{}
	var resTrimestres map[string]interface{}
	var planesPadre []map[string]interface{}
	var respuestaPost map[string]interface{}
	reporte := make(map[string]interface{})
	nuevo := true

	// Caso especial para el plan de acción, retomar avances de seguimiento de versiones anteriores
	if tipo == "61f236f525e40c582a0840d0" && plan["padre_plan_id"] != nil {
		nuevo = false
		var seguimientosPeticion []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+plan["dependencia_id"].(string)+",vigencia:"+plan["vigencia"].(string)+",formato:false,nombre:"+url.QueryEscape(plan["nombre"].(string)), &resPadres); err == nil {
			request.LimpiezaRespuestaRefactor(resPadres, &planesPadre)

			var seguimientosLlenos []map[string]interface{}
			var seguimientosVacios []map[string]interface{}

			for _, padre := range planesPadre {
				var resSeguimientos map[string]interface{}
				var seguimientos []map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+padre["_id"].(string), &resSeguimientos); err == nil {
					request.LimpiezaRespuestaRefactor(resSeguimientos, &seguimientos)
					seguimientosPeticion = seguimientos
					for _, seguimiento := range seguimientos {
						if (len(seguimientosLlenos) + len(seguimientosVacios)) <= 4 {
							if fmt.Sprintf("%v", seguimiento["dato"]) != "{}" {
								seguimientosLlenos = append(seguimientosLlenos, seguimiento)
							} else {
								seguimientosVacios = append(seguimientosVacios, seguimiento)
							}
						} else {
							break
						}
					}
				}
			}

			if len(seguimientosPeticion) == 0 && plan["nueva_estructura"].(bool) {
				nuevo = true
			}

			if !nuevo {
				var resActualizacion map[string]interface{}
				var resCreacion map[string]interface{}
				var resSeguimientoDetalle map[string]interface{}
				detalle := make(map[string]interface{})
				dato := make(map[string]interface{})
				var resEstado map[string]interface{}
				estado := map[string]interface{}{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				for _, seguimiento := range seguimientosVacios {
					// ? Inactivar el actual
					seguimiento["activo"] = false
					request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &resActualizacion, seguimiento)
					arrReportes = append(arrReportes, resActualizacion["Data"].(map[string]interface{}))
					// ? Crear el nuevo
					seguimiento["activo"] = true
					seguimiento["plan_id"] = plan_id
					delete(seguimiento, "_id")
					request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &resCreacion, seguimiento)
					arrReportes = append(arrReportes, resCreacion["Data"].(map[string]interface{}))
				}

				for _, seguimiento := range seguimientosLlenos {

					dato = map[string]interface{}{}
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)

					listAct := make([]string, 0, len(dato))
					for k := range dato {
						listAct = append(listAct, k)
					}
					for _, idxAct := range listAct {
						id, existe := dato[idxAct].(map[string]interface{})["id"].(string)
						if existe && id != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resSeguimientoDetalle); err == nil {
								request.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
								detalle = helpers.ConvertirStringJson(detalle)
								// ? Inactivar el actual
								detalle["activo"] = false
								helpers.GuardarDetalleSeguimiento(detalle, true) // true => PUT
								// ? crear el nuevo
								detalle["activo"] = true
								detalle["estado"] = estado
								delete(detalle, "_id")
								delete(detalle, "cuantitativo")
								newDetalleId := helpers.GuardarDetalleSeguimiento(detalle, false) // false => POST
								dato[idxAct].(map[string]interface{})["id"] = newDetalleId
							}
						}
					}

					// ? Inactiva el actual
					seguimiento["activo"] = false
					request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &resActualizacion, seguimiento)
					arrReportes = append(arrReportes, resActualizacion["Data"].(map[string]interface{}))

					// ? crear el nuevo
					seguimiento["activo"] = true
					seguimiento["plan_id"] = plan_id
					seguimiento["estado_seguimiento_id"] = "635c11e1e092c5fa5f099971" // En reporte
					valor, _ := json.Marshal(dato)
					str := string(valor)
					seguimiento["dato"] = str
					delete(seguimiento, "_id")
					request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &resCreacion, seguimiento)
					arrReportes = append(arrReportes, resCreacion["Data"].(map[string]interface{}))
				}
			}
		}
	}

	if nuevo {
		if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?query=Id:"+plan["dependencia_id"].(string), &resDependencia); err != nil {
			errRes = errors.New("error del servicio AvalarPlan: Error al consultar dependencia")
			respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
			if err != nil || respuesta["Success"] == false {
				errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
			}
			return nil, errRes
		}
		for i := 0; i < len(trimestres); i++ {
			periodo := int(trimestres[i]["Id"].(float64))
			if nuevaEstructura, ok := plan["nueva_estructura"].(bool); ok && nuevaEstructura {
				var respuestaRegistro map[string]interface{}
				planesInteresArray := []interface{}{
					map[string]interface{}{
						"_id":    plan["formato_id"],
						"nombre": plan["nombre"],
					},
				}
				planesInteresJSON, err := json.Marshal(planesInteresArray)
				if err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al convertir el array a JSON")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}

				dependenciaID, err := strconv.Atoi(plan["dependencia_id"].(string))
				if err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al convertir el ID de la dependencia")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}
				unidadInteresArray := []interface{}{
					map[string]interface{}{
						"Id":     dependenciaID,
						"Nombre": resDependencia[0]["Nombre"].(string),
					},
				}
				unidadInteresJSON, err := json.Marshal(unidadInteresArray)
				if err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al convertir el array a JSON")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}

				body := make(map[string]interface{})
				body["periodo_id"] = strconv.Itoa(periodo)
				body["planes_interes"] = string(planesInteresJSON)
				body["unidades_interes"] = string(unidadInteresJSON)
				body["tipo_seguimiento_id"] = tipo
				body["activo"] = true

				if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/1", "POST", &respuestaRegistro, body); err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al buscar el periodo-seguimiento")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}

				if respuestaRegistro["Data"] == nil {
					errRes = errors.New("error del servicio AvalarPlan: No se encontró el periodo-seguimiento")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}

				reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
				reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)
				if resDependencia[0]["Nombre"] != nil {
					reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + resDependencia[0]["Nombre"].(string)
				}
				reporte["activo"] = true
				reporte["plan_id"] = plan_id
				reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
				reporte["periodo_seguimiento_id"] = respuestaRegistro["Data"].([]interface{})[0].(map[string]interface{})["_id"]
				reporte["fecha_inicio"] = respuestaRegistro["Data"].([]interface{})[0].(map[string]interface{})["fecha_inicio"]
				reporte["tipo_seguimiento_id"] = tipo
				reporte["dato"] = "{}"

				_, err = json.Marshal(reporte)
				if err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al convertir el reporte a JSON")
					return nil, errRes
				}

				if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error creando reporte")
					return nil, errRes
				}

				arrReportes = append(arrReportes, respuestaPost["Data"].(map[string]interface{}))
				respuestaRegistro = nil
				respuestaPost = nil
			} else {
				// El parámetro 'nueva_estructura' no está presente o no es del tipo bool o no es true.
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:`+tipo+`,periodo_id:`+strconv.Itoa(periodo), &resTrimestres); err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error al consultar el periodo-seguimiento")
					respuesta, err := helpers.CambiarEstadoPlan(plan, id_estado_preaval)
					if err != nil || respuesta["Success"] == false {
						errRes = errors.New("error del servicio AvalarPlan: Error al cambiar estado del plan")
					}
					return nil, errRes
				}

				reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
				reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)
				if resDependencia[0]["Nombre"] != nil {
					reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + resDependencia[0]["Nombre"].(string)
				}
				reporte["activo"] = true
				reporte["plan_id"] = plan_id
				reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
				reporte["periodo_seguimiento_id"] = resTrimestres["Data"].([]interface{})[0].(map[string]interface{})["_id"]
				reporte["fecha_inicio"] = resTrimestres["Data"].([]interface{})[0].(map[string]interface{})["fecha_fin"]
				reporte["tipo_seguimiento_id"] = tipo
				reporte["dato"] = "{}"

				if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
					errRes = errors.New("error del servicio AvalarPlan: Error creando reporte")
					return nil, errRes
				}

				arrReportes = append(arrReportes, respuestaPost["Data"].(map[string]interface{}))
				respuestaPost = nil
			}
		}
	}
	return arrReportes, nil
}

func RevisarSeguimientoJefeDependencia(seguimiento_id string) (map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,_id:"+seguimiento_id, &respuesta); err != nil {
		return nil, errors.New("error del servicio RevisarSeguimientoJefeDependencia: No se pudo consultar el seguimiento")
	}

	aux := make([]map[string]interface{}, 1)
	request.LimpiezaRespuestaRefactor(respuesta, &aux)

	seguimiento = aux[0]
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	avalado, observacion, mensaje := seguimientoVerificable(seguimiento)

	if avalado || observacion {
		var codigo_abreviacion string

		if observacion {
			codigo_abreviacion = "RVCO"
		} else {
			codigo_abreviacion = "RV"
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
			estado := map[string]interface{}{
				"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
				"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
			}
			seguimiento["estado_seguimiento_id"] = estado["id"]
		}

		if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
			return nil, errors.New("error del servicio RevisarSeguimientoJefeDependencia: No se pudo actualizar el seguimiento")
		}

		data := respuesta["Data"].(map[string]interface{})
		return data, nil
	} else {
		responseJSON, _ := json.Marshal(mensaje)
		return nil, fmt.Errorf("%v", string(responseJSON))
	}
}

func seguimientoVerificable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false
	estado := map[string]interface{}{}

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		request.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

				actividades, errActividades := ConsultarActividades(subgrupos[i]["_id"].(string))
				if errActividades == nil {
					dato_plan_str := seguimiento["dato"].(string)
					json.Unmarshal([]byte(dato_plan_str), &datoPlan)

					for indexActividad, element := range datoPlan {
						id, segregado := element.(map[string]interface{})["id"]

						for _, actividad := range actividades {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								if indexActividad == actividad["index"] {
									if segregado && id != "" {
										if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
											request.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
											detalle = helpers.ConvertirStringJson(detalle)
											estado = detalle["estado"].(map[string]interface{})
											if estado["nombre"] != "Actividad Verificada" && estado["nombre"] != "Con observaciones" {
												dato[indexActividad] = actividad["dato"]
											}
										}
									} else {
										if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad Verificada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
											dato[indexActividad] = actividad["dato"]
										}
									}
								}
							} else if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
								if segregado && id != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
										request.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
										detalle = helpers.ConvertirStringJson(detalle)
										estado = detalle["estado"].(map[string]interface{})
										if estado["nombre"] != "Actividad Verificada" && estado["nombre"] != "Con observaciones" {
											dato[indexActividad] = actividad["dato"]
										}
									}
								} else {
									if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad Verificada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							}
						}

						if segregado && id != "" {
							if estado["nombre"] == "Con observaciones" {
								observaciones = true
							}

							if estado["nombre"] == "Actividad Verificada" {
								avaladas = true
							}
						} else {
							if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Con observaciones" {
								observaciones = true
							}

							if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Actividad Verificada" {
								avaladas = true
							}
						}
					}
				}
			}
		}
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return false, false, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}
	}

	return avaladas, observaciones, nil
}

func ObtenerPromedioBrechayEstado(requestBody []byte) (respuesta []map[string]interface{}, outputError error) {
	var body map[string]interface{}
	var trimestres []map[string]interface{}

	if err := json.Unmarshal(requestBody, &body); err == nil {
		nombrePlan := body["nombre"].(string)
		id := body["id"].(string)
		vigencia := body["vigencia"].(string)
		dependencia := body["dependencia"].(string)

		periodos, _ := ConsultarTrimestres(vigencia)
		if len(periodos) == 0 {
			outputError = errors.New("error del servicio ObtenerPromedioBrechayEstado: El plan no tiene definido Trimestres")
			return nil, outputError
		}

		for _, periodo := range periodos {
			trimestre := map[string]interface{}{
				"codigo": periodo["ParametroId"].(map[string]interface{})["CodigoAbreviacion"],
				"nombre": periodo["ParametroId"].(map[string]interface{})["Nombre"],
			}
			trimestres = append(trimestres, trimestre)
		}

		for _, tr := range trimestres {
			estado, err := ConsultarEstadoTrimestre(id, tr["codigo"].(string))
			if err != nil {
				logs.Error("Error --> ", err)
				return nil, errors.New(err.Error())
			}
			tr["estado"] = estado["estado_seguimiento_id"].(map[string]interface{})["nombre"]
		}

		unidades, errUnd := comunhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia)
		if errUnd != nil {
			logs.Error("Error --> ", errUnd)
			return nil, fmt.Errorf("error al obtener unidades: %v", errUnd)
		}
		if len(unidades) != 0 {
			planesPeriodo, errPer := comunhelper.GetPlanesPeriodo(dependencia, vigencia)
			if errPer != nil {
				logs.Error("Error --> ", errPer)
				return nil, fmt.Errorf("error al obtener periodos: %v", errPer)
			}

			var planEspecifico map[string]interface{}
			for _, item := range planesPeriodo {
				if item["plan"] == nombrePlan {
					planEspecifico = item
				}
			}
			if planEspecifico != nil {
				pers := planEspecifico["periodos"].([]map[string]interface{})
				ultimoPeriodo := pers[len(pers)-1]
				ultimoPeriodoID := ultimoPeriodo["id"].(string)

				periodosPlan := comunhelper.GetPeriodosPlan(vigencia, id)
				if len(periodosPlan) == 0 {
					outputError = errors.New("error del servicio ObtenerPromedioBrechayEstado: El plan no posee trimestres")
					return nil, outputError
				} else {
					var evaluacion []map[string]interface{}
					for posicionTri, tri := range periodosPlan {
						if tri["_id"] == ultimoPeriodoID {
							evaluacion = comunhelper.GetEvaluacion(planEspecifico["id"].(string), periodosPlan, posicionTri)
							break
						}
					}
					if evaluacion == nil {
						outputError = errors.New("error del servicio ObtenerPromedioBrechayEstado: El plan no posee evaluacion")
						return nil, outputError
					}
					var brechasT1 []float64
					var brechasT2 []float64
					var brechasT3 []float64
					var brechasT4 []float64
					for _, eval := range evaluacion {
						if eval["trimestre1"] != nil && eval["trimestre1"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre1"].(map[string]interface{})) > 0 {
							brechasT1 = append(brechasT1, eval["trimestre1"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre2"] != nil && eval["trimestre2"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre2"].(map[string]interface{})) > 0 {
							brechasT2 = append(brechasT2, eval["trimestre2"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre3"] != nil && eval["trimestre3"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre3"].(map[string]interface{})) > 0 {
							brechasT3 = append(brechasT3, eval["trimestre3"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre4"] != nil && eval["trimestre4"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre4"].(map[string]interface{})) > 0 {
							brechasT4 = append(brechasT4, eval["trimestre4"].(map[string]interface{})["brecha"].(float64))
						}
					}

					for _, tr := range trimestres {
						var suma float64 = 0
						var prod float64 = 0
						if tr["codigo"] == "T1" && len(brechasT1) != 0 {
							for _, numero := range brechasT1 {
								suma += numero
							}
							if len(brechasT1) > 0 {
								prod = (suma / float64(len(brechasT1)))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T2" && len(brechasT2) != 0 {
							for _, numero := range brechasT2 {
								suma += numero
							}
							if len(brechasT2) > 0 {
								prod = (suma / float64(len(brechasT2)))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T3" && len(brechasT3) != 0 {
							for _, numero := range brechasT3 {
								suma += numero
							}
							if len(brechasT3) > 0 {
								prod = (suma / float64(len(brechasT3)))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T4" && len(brechasT4) != 0 {
							for _, numero := range brechasT4 {
								suma += numero
							}
							if len(brechasT4) > 0 {
								prod = (suma / float64(len(brechasT4)))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else {
							tr["promedioBrechas"] = 0
						}
					}
				}

			} else {
				for _, tr := range trimestres {
					tr["promedioBrechas"] = 0
				}
			}
		} else {
			for _, tr := range trimestres {
				tr["promedioBrechas"] = 0
			}
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}

	return trimestres, outputError
}
