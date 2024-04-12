package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/planeacion"
	"github.com/udistrital/utils_oas/request"
)

func HabilitarReportes(entrada map[string]interface{}) ([]map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var respuestaPut map[string]interface{}
	var reportes []map[string]interface{}

	if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &respuesta); errPeriodo == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &reportes)

		if len(reportes) > 0 {
			var elemento = reportes[0]

			elemento["activo"] = true
			elemento["fecha_inicio"] = entrada["fecha_inicio"]
			elemento["fecha_fin"] = entrada["fecha_fin"]
			elemento["unidades_interes"] = "[]"
			if errPeriodoSeguimiento := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+elemento["_id"].(string), "PUT", &respuestaPut, elemento); errPeriodoSeguimiento != nil {
				logs.Error("Error -->", errPeriodoSeguimiento)
				return nil, errors.New(errPeriodoSeguimiento.Error())
			}
		} else {
			elemento := map[string]interface{}{
				"tipo_seguimiento_id": "61f236f525e40c582a0840d0",
				"activo":              true,
				"fecha_inicio":        entrada["fecha_inicio"],
				"fecha_fin":           entrada["fecha_fin"],
				"periodo_id":          entrada["periodo_id"],
				"unidades_interes":    "[]",
			}

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &respuestaPut, elemento); err == nil {
				if errPeriodoSeguimiento := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &respuesta); errPeriodoSeguimiento == nil {
					request.LimpiezaRespuestaRefactor(respuesta, &reportes)
					return reportes, nil
				} else {
					logs.Error("Error -->", errPeriodoSeguimiento)
					return nil, errors.New(errPeriodoSeguimiento.Error())
				}
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
		}
	} else {
		logs.Error("Error -->", errPeriodo)
		return nil, errors.New(errPeriodo.Error())
	}
	return nil, errors.New("error del servicio HabilitarReportes: No se habilitaron los reportes")
}

func CrearReportes(plan_identificador string, tipo string) ([]map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var respuestaPadres map[string]interface{}
	var respuestaDependencia []map[string]interface{}
	var respuestaTrimestres map[string]interface{}
	var plan map[string]interface{}
	var planesPadre []map[string]interface{}
	var respuestaPost map[string]interface{}
	var arregloReportes []map[string]interface{}
	reporte := make(map[string]interface{})
	nuevo := true

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_identificador, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &plan)
		trimestres, errTrimestres := ConsultarTrimestres(plan["vigencia"].(string))
		if errTrimestres != nil {
			logs.Error("Error -->", errTrimestres)
			return nil, errors.New(errTrimestres.Error())
		}

		// Caso especial para el plan de acción, retomar avances de seguimiento de versiones anteriores
		if tipo == "61f236f525e40c582a0840d0" && plan["padre_plan_id"] != nil {
			nuevo = false

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+plan["dependencia_id"].(string)+",vigencia:"+plan["vigencia"].(string)+",formato:false,nombre:"+url.QueryEscape(plan["nombre"].(string)), &respuestaPadres); err == nil {
				request.LimpiezaRespuestaRefactor(respuestaPadres, &planesPadre)

				var seguimientosLlenos []map[string]interface{}
				var seguimientosVacios []map[string]interface{}

				for _, padre := range planesPadre {
					var respuestaSeguimientos map[string]interface{}
					var seguimientos []map[string]interface{}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+padre["_id"].(string), &respuestaSeguimientos); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientos, &seguimientos)

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
					} else {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}
				}

				var respuestaActualizacion map[string]interface{}
				var respuestaCreacion map[string]interface{}
				var respuestaSeguimientoDetalle map[string]interface{}
				detalle := make(map[string]interface{})
				dato := make(map[string]interface{})
				var respuestaEstado map[string]interface{}
				estado := map[string]interface{}{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				} else {
					logs.Error("Error -->", err)
					return nil, errors.New(err.Error())
				}

				for _, seguimiento := range seguimientosVacios {
					// Inactivar el actual
					seguimiento["activo"] = false
					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuestaActualizacion, seguimiento); err == nil {
						arregloReportes = append(arregloReportes, respuestaActualizacion["Data"].(map[string]interface{}))
						// Crear el nuevo
						seguimiento["activo"] = true
						seguimiento["plan_id"] = plan_identificador
						delete(seguimiento, "_id")
						if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaCreacion, seguimiento); err == nil {
							arregloReportes = append(arregloReportes, respuestaCreacion["Data"].(map[string]interface{}))
						} else {
							logs.Error("Error -->", err)
							return nil, errors.New(err.Error())
						}
					} else {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}
				}

				for _, seguimiento := range seguimientosLlenos {
					dato = map[string]interface{}{}
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)
					listaActividades := make([]string, 0, len(dato))

					for k := range dato {
						listaActividades = append(listaActividades, k)
					}

					for _, indexActividad := range listaActividades {
						identificador, existe := dato[indexActividad].(map[string]interface{})["id"].(string)

						if existe && identificador != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+identificador, &respuestaSeguimientoDetalle); err == nil {
								request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
								planeacion.ConvertirStringJson(detalle)
								// Inactivar el actual
								detalle["activo"] = false
								guardarDetalleSeguimiento(detalle, true) // true => PUT
								// Crear el nuevo
								detalle["activo"] = true
								detalle["estado"] = estado
								delete(detalle, "_id")
								delete(detalle, "cuantitativo")
								NuevoDetalleIdentificador := guardarDetalleSeguimiento(detalle, false) // false => POST
								dato[indexActividad].(map[string]interface{})["id"] = NuevoDetalleIdentificador
							}
						}
					}
					// Inactiva el actual
					seguimiento["activo"] = false

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuestaActualizacion, seguimiento); err == nil {
						arregloReportes = append(arregloReportes, respuestaActualizacion["Data"].(map[string]interface{}))
					} else {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}

					// Crear el nuevo
					seguimiento["activo"] = true
					seguimiento["plan_id"] = plan_identificador
					seguimiento["estado_seguimiento_id"] = "635c11e1e092c5fa5f099971" // En reporte
					valor, _ := json.Marshal(dato)
					str := string(valor)
					seguimiento["dato"] = str
					delete(seguimiento, "_id")

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaCreacion, seguimiento); err == nil {
						arregloReportes = append(arregloReportes, respuestaCreacion["Data"].(map[string]interface{}))
					} else {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}
				}
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New(err.Error())
			}
		}

		if nuevo {
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?query=Id:"+plan["dependencia_id"].(string), &respuestaDependencia); err == nil {
			} else {
				logs.Error("Error -->", err)
				return nil, errors.New("error del servicio CrearReportes:    Error obteniendo la dependencia" + err.Error())
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
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error convirtiendo el array a JSON" + err.Error())
					}

					dependenciaID, err := strconv.Atoi(plan["dependencia_id"].(string))
					if err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error convirtiendo el ID de la dependencia a entero" + err.Error())
					}
					unidadInteresArray := []interface{}{
						map[string]interface{}{
							"Id":     dependenciaID,
							"Nombre": respuestaDependencia[0]["Nombre"].(string),
						},
					}
					unidadInteresJSON, err := json.Marshal(unidadInteresArray)
					if err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error convirtiendo el array a JSON" + err.Error())
					}

					body := make(map[string]interface{})
					body["periodo_id"] = strconv.Itoa(periodo)
					body["planes_interes"] = string(planesInteresJSON)
					body["unidades_interes"] = string(unidadInteresJSON)
					body["tipo_seguimiento_id"] = tipo
					body["activo"] = true

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/1", "POST", &respuestaRegistro, body); err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error buscando periodo-seguimiento" + err.Error())
					}

					reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
					reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)
					if respuestaDependencia[0]["Nombre"] != nil {
						reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + respuestaDependencia[0]["Nombre"].(string)
					}
					reporte["activo"] = true
					reporte["plan_id"] = plan_identificador
					reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
					reporte["periodo_seguimiento_id"] = respuestaRegistro["Data"].([]interface{})[0].(map[string]interface{})["_id"]
					reporte["fecha_inicio"] = respuestaRegistro["Data"].([]interface{})[0].(map[string]interface{})["fecha_inicio"]
					reporte["tipo_seguimiento_id"] = tipo
					reporte["dato"] = "{}"

					if _, err := json.Marshal(reporte); err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error al convertir el reporte a JSON" + err.Error())
					}

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New("error del servicio CrearReportes: Error creando reporte" + err.Error())
					}

					arregloReportes = append(arregloReportes, respuestaPost["Data"].(map[string]interface{}))
					respuestaRegistro = nil
					respuestaPost = nil

				} else {
					// El parámetro 'nueva_estructura' no está presente o no es del tipo bool o no es true.
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:`+tipo+`,periodo_id:`+strconv.Itoa(periodo), &respuestaTrimestres); err == nil {
						reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
						reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)
						if respuestaDependencia[0]["Nombre"] != nil {
							reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + respuestaDependencia[0]["Nombre"].(string)
						}
						reporte["activo"] = true
						reporte["plan_id"] = plan_identificador
						reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
						reporte["periodo_seguimiento_id"] = respuestaTrimestres["Data"].([]interface{})[0].(map[string]interface{})["_id"]
						reporte["fecha_inicio"] = respuestaTrimestres["Data"].([]interface{})[0].(map[string]interface{})["fecha_fin"]
						reporte["tipo_seguimiento_id"] = tipo
						reporte["dato"] = "{}"

						if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
							logs.Error("Error -->", err)
							return nil, errors.New("error del servicio CrearReportes: Error creando reporte" + err.Error())
						}

						arregloReportes = append(arregloReportes, respuestaPost["Data"].(map[string]interface{}))
						respuestaPost = nil
					} else {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}
				}
			}
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}

	return arregloReportes, nil
}

func ReportarSeguimiento(identificadorSeguimiento string) (interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+identificadorSeguimiento, &respuesta); err == nil {
		aux := make(map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux
		reportable, mensaje := seguimientoReportable(seguimiento)

		if reportable {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:EAR", &respuestaEstado); err == nil {
				seguimiento["estado_seguimiento_id"] = respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
			}

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				logs.Error("Error -->", err)
				return nil, errors.New("error del servicio ReportarSeguimiento:    Error actualizando seguimiento" + err.Error())
			}

			return seguimiento, nil
		} else {
			logs.Error("Error -->", formatdata.JsonPrint(mensaje))
			return nil, errors.New("error del servicio ReportarSeguimiento:    El seguimiento no es reportable")
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func ReportarActividad(requestBody []byte, indiceActividad string) (interface{}, error) {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var estado map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})

	if err := json.Unmarshal(requestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+body["SeguimientoId"].(string), &respuesta); err == nil {
			aux := make(map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)
			reportable, mensaje := actividadReportable(seguimiento, indiceActividad)

			if reportable {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				identificador, segregable := dato[indiceActividad].(map[string]interface{})["id"].(string)

				if segregable && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = planeacion.ConvertirStringJson(detalle)
						detalle["estado"] = estado
						guardarDetalleSeguimiento(detalle, true)
					}
				} else {
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
					b, _ := json.Marshal(dato)
					str := string(b)
					seguimiento["dato"] = str

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						logs.Error("Error -->", err)
						return nil, errors.New(err.Error())
					}
				}
				return seguimiento, nil
			} else {
				logs.Error("Error -->", formatdata.JsonPrint(mensaje))
				return nil, errors.New("error del servicio ReportarActividad:    La actividad no es reportable")
			}
		} else {
			logs.Error("Error -->", err)
			return nil, errors.New(err.Error())
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}
