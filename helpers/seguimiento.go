package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/planeacion"
	"github.com/udistrital/utils_oas/request"
)

func ConsultarTrimestres(vigencia string) []map[string]interface{} {
	var respuesta map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}
	var respuestaParametros map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro?query=CodigoAbreviacion.in:T1|T2|T3|T4", &respuestaParametros); err == nil {
		var parametros []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaParametros, &parametros)

		for _, parametro := range parametros {
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:"+parametro["CodigoAbreviacion"].(string), &respuesta); err == nil {
				request.LimpiezaRespuestaRefactor(respuesta, &trimestre)
				trimestres = append(trimestres, trimestre...)
			} else {
				panic(map[string]interface{}{"funcion": "ConsultarTrimestres", "err": "Error ", "status": "400", "log": err})
			}
			trimestre = nil
		}
	} else {
		panic(map[string]interface{}{"funcion": "ConsultarTrimestres", "err": "Error ", "status": "400", "log": err})
	}

	return trimestres
}

func GuardarDetalleSeguimiento(detalle map[string]interface{}, actualizar bool) string {
	var respuesta map[string]interface{}
	var identificador string

	detalle = planeacion.ConvertirJsonString(detalle)

	if _, existe := detalle["informacion"]; !existe {
		detalle["informacion"] = "{}"
	}
	if _, existe := detalle["cualitativo"]; !existe {
		detalle["cualitativo"] = "{}"
	}
	if _, existe := detalle["cuantitativo"]; !existe {
		detalle["cuantitativo"] = "{}"
	}
	if _, existe := detalle["evidencia"]; !existe {
		detalle["evidencia"] = "[]"
	}

	if actualizar {
		if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+detalle["_id"].(string), "PUT", &respuesta, detalle); err == nil {
			aux := make(map[string]interface{})
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			identificador = aux["_id"].(string)
		}
	} else {
		if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle", "POST", &respuesta, detalle); err == nil {
			aux := make(map[string]interface{})
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			identificador = aux["_id"].(string)
		}
	}
	return identificador
}

func ConsultarActividades(subgrupo_identificador string) []map[string]interface{} {
	var respuesta map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var actividades []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo_identificador, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		subgrupoDetalle = aux[0]

		if subgrupoDetalle["dato_plan"] != nil {
			dato_plan_str := subgrupoDetalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &datoPlan)

			for indiceActividad, elemento := range datoPlan {
				_ = indiceActividad
				if err != nil {
					log.Panic(err)
				}
				if elemento.(map[string]interface{})["activo"] == true {
					actividades = append(actividades, elemento.(map[string]interface{}))
				}

				request.SortSlice(&actividades, "index")
			}
		}
	} else {
		panic(map[string]interface{}{"Code": "400", "Body": err, "Type": "error"})
	}
	return actividades
}

func ConsultarEstadoSeguimiento(seguimiento map[string]interface{}) string {
	var respuestaEstado map[string]interface{}
	enReporte := true
	estado := map[string]interface{}{}
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
					panic(err)
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
		panic(err)
	}
	return estado["id"].(string)
}

func ConsultarActividad(seguimiento map[string]interface{}, indice string, trimestre string) map[string]interface{} {
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
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indice].(map[string]interface{})["id"].(string), &respuestaDetalle); err == nil {
				if respuestaDetalle["Data"] != "null" {
					request.LimpiezaRespuestaRefactor(respuestaDetalle, &detalle)
					detalle = planeacion.ConvertirStringJson(detalle)
					identificador = detalle["_id"].(string)

					if len(detalle["informacion"].(map[string]interface{})) == 0 {
						informacion = ConsultarInformacionPlan(seguimiento, indice)
					} else {
						informacion = detalle["informacion"].(map[string]interface{})
					}

					if len(detalle["cuantitativo"].(map[string]interface{})) == 0 {
						cuantitativo = ConsultarCuantitativoPlan(seguimiento, indice, trimestre)
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
				informacion = ConsultarInformacionPlan(seguimiento, indice)
			} else {
				informacion = dato[indice].(map[string]interface{})["informacion"].(map[string]interface{})
			}

			if dato[indice].(map[string]interface{})["cuantitativo"] == nil {
				cuantitativo = ConsultarCuantitativoPlan(seguimiento, indice, trimestre)
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
		informacion = ConsultarInformacionPlan(seguimiento, indice)
		cuantitativo = ConsultarCuantitativoPlan(seguimiento, indice, trimestre)
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

func ConsultarInformacionPlan(seguimiento map[string]interface{}, indice string) map[string]interface{} {
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
		panic(err)
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &respuestaPeriodoSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(respuestaPeriodoSeguimiento, &periodoSeguimiento)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &respuestaPeriodo); err == nil {
			request.LimpiezaRespuestaRefactor(respuestaPeriodo, &periodo)
			informacion["trimestre"] = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
		} else {
			panic(err)
		}
	} else {
		panic(err)
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
					panic(err)
				}
			}
		}
	} else {
		panic(err)
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+informacion["unidad"].(string), &respuestaDependencia); err == nil {
		informacion["unidad"] = respuestaDependencia[0]["DependenciaId"].(map[string]interface{})["Nombre"]
	} else {
		informacion["unidad"] = nil
	}
	return informacion
}

func ConsultarCuantitativoPlan(seguimiento map[string]interface{}, indice string, trimestre string) map[string]interface{} {
	var respuestaInformacion map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var hijos []interface{}
	var subgrupos []map[string]interface{}
	var indicadores []map[string]interface{}
	var respuestas []map[string]interface{}
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
								panic(err)
							}
						}

						if informacion["reporteDenominador"] == 1.0 {
							informacion["reporteDenominador"] = nil
						}

						if informacion["nombre"] != nil && informacion["nombre"] != "" {
							indicadores = append(indicadores, informacion)
							respuestas = append(respuestas, respuesta)
						}

						respuestas = ConsultarRespuestaAnterior(seguimiento, len(indicadores)-1, respuestas, indice, trimestre)
					} else {
						panic(err)
					}
				}
				break
			}
		}
	} else {
		panic(err)
	}

	response["indicadores"] = indicadores
	response["resultados"] = respuestas
	return response
}

func ConsultarRespuestaAnterior(dataSeg map[string]interface{}, indice int, respuestas []map[string]interface{}, indiceActividad string, trimestre string) []map[string]interface{} {
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
									detalle = planeacion.ConvertirStringJson(detalle)

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
									panic(err)
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
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	}
	return respuestas
}

func ActividadConObservaciones(seguimiento map[string]interface{}) bool {
	var cuantitativo map[string]interface{}
	var cualitativo map[string]interface{}

	if seguimiento["cuantitativo"] != nil {
		cuantitativo = seguimiento["cuantitativo"].(map[string]interface{})
		for _, indicador := range cuantitativo["indicadores"].([]interface{}) {
			if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" && indicador.(map[string]interface{})["observaciones"] != nil {
				return true
			}
		}
	}

	if seguimiento["cualitativo"] != nil {
		cualitativo = seguimiento["cualitativo"].(map[string]interface{})
		if cualitativo["observaciones"] != "" && cualitativo["observaciones"] != "Sin observación" && cualitativo["observaciones"] != nil {
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

func ActividadReportable(seguimiento map[string]interface{}, indiceActividad string) (bool, map[string]interface{}) {
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	var cuantitativo interface{}
	var cualitativo interface{}
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[indiceActividad] == nil {
		return false, map[string]interface{}{"error": 1, "motivo": "Actividad sin seguimiento"}
	} else {
		_, datosUnidos := dato[indiceActividad].(map[string]interface{})["estado"]

		if datosUnidos {
			estado = dato[indiceActividad].(map[string]interface{})["estado"].(map[string]interface{})
			cuantitativo = dato[indiceActividad].(map[string]interface{})["cuantitativo"]
			cualitativo = dato[indiceActividad].(map[string]interface{})["cualitativo"]
		} else {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
				request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
				detalle = planeacion.ConvertirStringJson(detalle)
				estado = detalle["estado"].(map[string]interface{})
				cualitativo = detalle["cualitativo"]
				cuantitativo = detalle["cuantitativo"]
			}

			if fmt.Sprintf("%v", cuantitativo) == "map[]" {
				cuantitativo = nil
			}

			if fmt.Sprintf("%v", cualitativo) == "map[]" {
				cualitativo = nil
			}
		}

		if estado["nombre"] != "Actividad en reporte" {
			return false, map[string]interface{}{"error": 2, "motivo": "El estado de la actividad no es el adecuado"}
		}

		if cuantitativo == nil {
			return false, map[string]interface{}{"error": 3, "motivo": "Componenten cuantitativo sin guardar"}
		}

		if cualitativo == nil {
			return false, map[string]interface{}{"error": 4, "motivo": "Componenten cualitativo sin guardar"}
		} else {
			cualitativo := cualitativo.(map[string]interface{})
			if cualitativo["dificultades"] == "" || cualitativo["productos"] == "" || cualitativo["reporte"] == "" {
				return false, map[string]interface{}{"error": 5, "motivo": "Campos vacios en el componenten cualitativo"}
			}
		}
	}
	return true, nil
}

func SeguimientoReportable(seguimiento map[string]interface{}) (bool, map[string]interface{}) {
	var respuesta map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	planIdentificador := seguimiento["plan_id"].(string)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planIdentificador, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
				actividades := ConsultarActividades(subgrupos[i]["_id"].(string))

				if seguimiento["dato"] == "{}" {
					for _, actividad := range actividades {
						dato[actividad["index"].(string)] = actividad["dato"]
					}
					return false, map[string]interface{}{"error": 1, "motivo": "No hay actividades resportadas", "actividades": dato}
				} else {
					dato_plan_str := seguimiento["dato"].(string)
					json.Unmarshal([]byte(dato_plan_str), &datoPlan)

					for indiceActividad, elemento := range datoPlan {
						identificador, segregado := elemento.(map[string]interface{})["id"]

						if segregado && identificador != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+identificador.(string), &respuestaSeguimientoDetalle); err == nil {
								request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
								detalle = planeacion.ConvertirStringJson(detalle)
							} else {
								panic(err)
							}

							for _, actividad := range actividades {
								if reflect.TypeOf(actividad["index"]).String() == "string" {
									if indiceActividad == actividad["index"] {
										actividad["estado"] = detalle["estado"]
									}
								} else {
									if indiceActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
										actividad["estado"] = detalle["estado"]
									}
								}
							}
						} else {
							for _, actividad := range actividades {
								if reflect.TypeOf(actividad["index"]).String() == "string" {
									if indiceActividad == actividad["index"] {
										actividad["estado"] = elemento.(map[string]interface{})["estado"]
									}
								} else {
									if indiceActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
										actividad["estado"] = elemento.(map[string]interface{})["estado"]
									}
								}
							}
						}
					}
					for _, actividad := range actividades {
						if actividad["estado"] == nil {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						} else if actividad["estado"].(map[string]interface{})["nombre"] != "Actividad reportada" && actividad["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						}
					}

					if fmt.Sprintf("%v", dato) != "map[]" {
						return false, map[string]interface{}{"error": 2, "motivo": "Hay actividades sin resportar", "actividades": dato}
					} else {
						return true, nil
					}
				}
			}
		}
	} else {
		panic(err)
	}
	return true, nil
}

func SeguimientoAvalable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}) {
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
				actividades := ConsultarActividades(subgrupos[i]["_id"].(string))
				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				for indiceActividad, elemento := range datoPlan {
					identificador, segregado := elemento.(map[string]interface{})["id"]

					for _, actividad := range actividades {
						if reflect.TypeOf(actividad["index"]).String() == "string" {
							if indiceActividad == actividad["index"] {
								if segregado && identificador != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+elemento.(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
										request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
										detalle = planeacion.ConvertirStringJson(detalle)
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
									detalle = planeacion.ConvertirStringJson(detalle)
									estado = detalle["estado"].(map[string]interface{})
									if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
										dato[indiceActividad] = actividad["dato"]
									}
								} else {
									panic(err)
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
	} else {
		panic(err)
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return false, false, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}
	}

	return avaladas, observaciones, nil
}
