package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/helpers"
	"github.com/udistrital/utils_oas/planeacion"
	"github.com/udistrital/utils_oas/request"
)

// SeguimientoController operations for Seguimiento
type SeguimientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SeguimientoController) URLMapping() {
	c.Mapping("HabilitarReportes", c.HabilitarReportes)
	c.Mapping("CrearReportes", c.CrearReportes)
	c.Mapping("ConsultarPeriodos", c.ConsultarPeriodos)
	c.Mapping("ConsultarActividadesGenerales", c.ConsultarActividadesGenerales)
	c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	c.Mapping("ConsultarSeguimiento", c.ConsultarSeguimiento)
	c.Mapping("ConsultarIndicadores", c.ConsultarIndicadores)
	c.Mapping("ConsultarAvanceIndicador", c.ConsultarAvanceIndicador)
	c.Mapping("ConsultarEstadoTrimestre", c.ConsultarEstadoTrimestre)
	c.Mapping("GuardarDocumentos", c.GuardarDocumentos)
	c.Mapping("GuardarCualitativo", c.GuardarCualitativo)
	c.Mapping("GuardarCuantitativo", c.GuardarCuantitativo)
	c.Mapping("ReportarActividad", c.ReportarActividad)
	c.Mapping("ReportarSeguimiento", c.ReportarSeguimiento)
	c.Mapping("RetornarActividad", c.RetornarActividad)
	c.Mapping("MigrarInformacion", c.MigrarInformacion)
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description put Seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403
// @router /habilitar_reportes [put]
func (c *SeguimientoController) HabilitarReportes() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var respuesta map[string]interface{}
	var respuestaPut map[string]interface{}
	var reportes []map[string]interface{}
	var entrada map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &reportes)

		if len(reportes) > 0 {
			var elemento = reportes[0]

			elemento["activo"] = true
			elemento["fecha_inicio"] = entrada["fecha_inicio"]
			elemento["fecha_fin"] = entrada["fecha_fin"]
			elemento["unidades_interes"] = "[]"
			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+elemento["_id"].(string), "PUT", &respuestaPut, elemento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
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
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &respuesta); err == nil {
					request.LimpiezaRespuestaRefactor(respuesta, &reportes)
					c.Data["json"] = reportes
				} else {
					panic(err)
				}
			} else {
				panic(map[string]interface{}{"funcion": "HabilitarReportes", "err": "Error actualizando periodo-seguimiento", "status": "400", "log": err})
			}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// CrearReportes ...
// @Title CrearReportes
// @Description Post Seguimiento
// @Param	plan 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /crear_reportes/:plan/:tipo [post]
func (c *SeguimientoController) CrearReportes() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "InversionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	plan_identificador := c.Ctx.Input.Param(":plan")
	tipo := c.Ctx.Input.Param(":tipo")
	var respuesta map[string]interface{}
	var respuestaPadres map[string]interface{}
	var respuestaDependencia map[string]interface{}
	var respuestaTrimestres map[string]interface{}
	var plan map[string]interface{}
	var planesPadre []map[string]interface{}
	var respuestaPost map[string]interface{}
	var arregloReportes []map[string]interface{}
	reporte := make(map[string]interface{})
	nuevo := true

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_identificador, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &plan)
		trimestres := helpers.ConsultarTrimestres(plan["vigencia"].(string))

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
						panic(err)
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
					panic(err)
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
							panic(err)
						}
					} else {
						panic(err)
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
								helpers.GuardarDetalleSeguimiento(detalle, true) // true => PUT
								// Crear el nuevo
								detalle["activo"] = true
								detalle["estado"] = estado
								delete(detalle, "_id")
								delete(detalle, "cuantitativo")
								NuevoDetalleIdentificador := helpers.GuardarDetalleSeguimiento(detalle, false) // false => POST
								dato[indexActividad].(map[string]interface{})["id"] = NuevoDetalleIdentificador
							}
						}
					}
					// Inactiva el actual
					seguimiento["activo"] = false

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuestaActualizacion, seguimiento); err == nil {
						arregloReportes = append(arregloReportes, respuestaActualizacion["Data"].(map[string]interface{}))
					} else {
						panic(err)
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
						panic(err)
					}
				}
			} else {
				panic(err)
			}
		}

		if nuevo {
			for i := 0; i < len(trimestres); i++ {
				periodo := int(trimestres[i]["Id"].(float64))

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:`+tipo+`,periodo_id:`+strconv.Itoa(periodo), &respuestaTrimestres); err == nil {
					reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
					reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)

					if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"dependencia/"+plan["dependencia_id"].(string), &respuestaDependencia); err == nil {
						if respuestaDependencia["Nombre"] != nil {
							reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + respuestaDependencia["Nombre"].(string)
						}
					} else {
						panic(err)
					}

					reporte["activo"] = false
					reporte["plan_id"] = plan_identificador
					reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
					reporte["periodo_seguimiento_id"] = respuestaTrimestres["Data"].([]interface{})[0].(map[string]interface{})["_id"]
					reporte["fecha_inicio"] = respuestaTrimestres["Data"].([]interface{})[0].(map[string]interface{})["fecha_fin"]
					reporte["tipo_seguimiento_id"] = tipo
					reporte["dato"] = "{}"

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
						panic(map[string]interface{}{"funcion": "CrearReportes", "err": "Error creando reporte", "status": "400", "log": err})
					}

					arregloReportes = append(arregloReportes, respuestaPost["Data"].(map[string]interface{}))
					respuestaPost = nil
				} else {
					panic(err)
				}
			}
		}
	} else {
		panic(err)
	}
	c.Data["json"] = arregloReportes
	c.ServeJSON()
}

// ConsultarPeriodos ...
// @Title ConsultarPeriodos
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /consultar_periodos/:vigencia [get]
func (c *SeguimientoController) ConsultarPeriodos() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	vigencia := c.Ctx.Input.Param(":vigencia")

	if len(vigencia) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}

	trimestres := helpers.ConsultarTrimestres(vigencia)

	if len(trimestres) == 0 || trimestres[0]["Id"] == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}

	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": trimestres}
	}
	c.ServeJSON()
}

// ConsultarActividadesGenerales ...
// @Title ConsultarActividadesGenerales
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_actividades/:seguimiento_id [get]
func (c *SeguimientoController) ConsultarActividadesGenerales() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	seguimiento_identificador := c.Ctx.Input.Param(":seguimiento_id")
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
						actividades := helpers.ConsultarActividades(subgrupos[i]["_id"].(string))

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
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": actividades}
						break
					}
				}
			}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /guardar_seguimiento/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarSeguimiento() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var respuestaEstado map[string]interface{}
	var estadoSeguimiento string
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
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
				panic(err)
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
						detalle = planeacion.ConvertirStringJson(detalle)
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
			estadoSeguimiento = helpers.ConsultarEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarSeguimiento", "err": "Error actualizando seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// ConsultarSeguimiento ...
// @Title ConsultarSeguimiento
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_seguimiento/:plan_id/:index/:trimestre [get]
func (c *SeguimientoController) ConsultarSeguimiento() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestreIdentificador := c.Ctx.Input.Param(":trimestre")
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
				panic(err)
			}
		} else {
			panic(err)
		}

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)
		actividad, _ := json.Marshal(helpers.ConsultarActividad(seguimiento, indiceActividad, trimestre))
		json.Unmarshal([]byte(string(actividad)), &seguimientoActividad)
		seguimientoActividad["_id"] = seguimiento["_id"].(string)

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &respuestaEstado); err == nil {
			seguimientoActividad["estadoSeguimiento"] = respuestaEstado["Data"].(map[string]interface{})["nombre"].(string)
		} else {
			panic(err)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimientoActividad}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// ConsultarIndicadores ...
// @Title ConsultarIndicadores
// @Description get Seguimiento
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_indicadores/:plan_id [get]
func (c *SeguimientoController) ConsultarIndicadores() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	plan_identificador := c.Ctx.Input.Param(":plan_id")
	var respuesta map[string]interface{}
	var subgrupos []map[string]interface{}
	var hijos []map[string]interface{}
	var indicadores []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_identificador, &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "indicador") {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+subgrupos[i]["_id"].(string), &respuesta); err == nil {
					request.LimpiezaRespuestaRefactor(respuesta, &hijos)

					for j := range hijos {
						if strings.Contains(strings.ToLower(hijos[j]["nombre"].(string)), "indicador") {
							aux := hijos[j]
							indicadores = append(indicadores, aux)
						}
					}

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": indicadores}
				} else {
					panic(err)
				}
				break
			}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// ConsultarAvanceIndicador ...
// @Title ConsultarAvanceIndicador
// @Description post Seguimiento by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /consultar_avance [post]
func (c *SeguimientoController) ConsultarAvanceIndicador() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var avancedata []map[string]interface{}
	var respuesta1 map[string]interface{}
	var avancedata1 []map[string]interface{}
	var respuesta2 map[string]interface{}
	var respuestaName map[string]interface{}
	var parametro_periodo_name []map[string]interface{}
	var avancedata2 []map[string]interface{}
	var parametro_periodo []map[string]interface{}
	var dato map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimiento1 map[string]interface{}
	var test1 string
	var periodoIdentificadorString string
	var periodoIdentificador float64
	var avanceAcumulado string
	var testavancePeriodo string
	var nombrePeriodo string
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+body["periodo_seguimiento_id"].(string), &respuesta); err == nil {
		request.LimpiezaRespuestaRefactor(respuesta, &avancedata)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_seguimiento_id"].(string), &respuesta); err == nil {
			request.LimpiezaRespuestaRefactor(respuesta, &parametro_periodo)
			parametroIdentificadorlen := parametro_periodo[0]
			parametroIdentificador := parametroIdentificadorlen["ParametroId"].(map[string]interface{})

			if parametroIdentificador["CodigoAbreviacion"] != "T1" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+body["periodo_seguimiento_id"].(string), &respuesta1); err == nil {
					request.LimpiezaRespuestaRefactor(respuesta1, &avancedata1)
					seguimiento1 = avancedata1[0]
					datoStrUltimoTrimestre := seguimiento1["dato"].(string)

					if datoStrUltimoTrimestre == "{}" {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 32)
						if err != nil {
							fmt.Println(err)
						}
						periodoIdentificador = priodoId_rest - 1
					} else {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 32)
						if err != nil {
							fmt.Println(err)
						}
						periodoIdentificador = priodoId_rest
					}
				} else {
					panic(err)
				}
				periodoIdentificadorString = fmt.Sprint(periodoIdentificador)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+periodoIdentificadorString, &respuesta2); err == nil {
					request.LimpiezaRespuestaRefactor(respuesta2, &avancedata2)
					seguimiento = avancedata2[0]
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)
					indicador1 := dato[body["index"].(string)].(map[string]interface{})
					avanceIndicador1 := indicador1[body["Nombre_del_indicador"].(string)].(map[string]interface{})
					avanceAcumulado = avanceIndicador1["avanceAcumulado"].(string)
					testavancePeriodo = avanceIndicador1["avancePeriodo"].(string)
				} else {
					panic(err)
				}

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_seguimiento_id"].(string), &respuestaName); err == nil {
					request.LimpiezaRespuestaRefactor(respuestaName, &parametro_periodo_name)
					paramIdlenName := parametro_periodo_name[0]

					paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
					nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
				} else {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
		avancePeriodo := body["avancePeriodo"].(string)
		aPe, err := strconv.ParseFloat(avancePeriodo, 32)
		if err != nil {
			fmt.Println(err)
		}

		aAc, err := strconv.ParseFloat(avanceAcumulado, 32)
		if err != nil {
			fmt.Println(err)
		}
		totalAcumulado := fmt.Sprint(aPe + aAc)
		generalData := make(map[string]interface{})
		generalData["avancePeriodo"] = avancePeriodo
		generalData["periodIdString"] = periodoIdentificadorString
		generalData["avanceAcumulado"] = totalAcumulado
		generalData["avancePeriodoPrev"] = testavancePeriodo
		generalData["avanceAcumuladoPrev"] = avanceAcumulado
		generalData["nombrePeriodo"] = nombrePeriodo
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": generalData}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// ConsultarEstadoTrimestre ...
// @Title ConsultarEstadoTrimestre
// @Description get Seguimiento del trimestre correspondiente
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404 not found resource
// @router /consultar_estado_trimestre/:plan_id/:trimestre [get]
func (c *SeguimientoController) ConsultarEstadoTrimestre() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var respuestaSeguimiento map[string]interface{}
	var respuestaPeriodoSeguimiento map[string]interface{}
	var respuestaPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var periodoSeguimiento []map[string]interface{}

	planIdentificador := c.Ctx.Input.Param(":plan_id")
	trimestre := c.Ctx.Input.Param(":trimestre")

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

									c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": plan}
									break
								}
							}
						}
					}
				}
			}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GuardarDocumentos ...
// @Title GuardarDocumentos
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_documentos/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarDocumentos() {
	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var estadoSeguimiento string
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
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
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else if comentario {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:CO", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
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
						detalle = planeacion.ConvertirStringJson(detalle)
						detalle["evidencia"] = evidencias
						detalle["estado"] = estado
						helpers.GuardarDetalleSeguimiento(detalle, true)
					} else {
						panic(err)
					}
				} else {
					dato[indiceActividad].(map[string]interface{})["evidencia"] = evidencias
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
				}
			}
			valor, _ := json.Marshal(dato)
			str := string(valor)
			seguimiento["dato"] = str
			estadoSeguimiento = helpers.ConsultarEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarDocumentos", "err": "Error guardado documentos del seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"seguimiento": detalle["evidencia"], "estadoActividad": estado}}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	}
	c.ServeJSON()
}

// GuardarCualitativo ...
// @Title GuardarCualitativo
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cualitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCualitativo() {
	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cualitativo map[string]interface{}
	var informacion map[string]interface{}
	var estadoSeguimiento string
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
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
					panic(err)
				}

				detalle = map[string]interface{}{"estado": estado, "cualitativo": cualitativo, "informacion": informacion}
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(detalle, false)}
			} else {
				identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]

				if segregado && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = planeacion.ConvertirStringJson(detalle)
						estado = detalle["estado"].(map[string]interface{})
					} else {
						panic(err)
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

					observacion = helpers.ActividadConObservaciones(body)
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
						panic(err)
					}
				}

				if segregado && identificador != "" {
					detalle["cualitativo"] = cualitativo
					detalle["informacion"] = informacion
					detalle["estado"] = estado
					helpers.GuardarDetalleSeguimiento(detalle, true)
				} else {
					dato[indiceActividad].(map[string]interface{})["informacion"] = informacion
					dato[indiceActividad].(map[string]interface{})["cualitativo"] = cualitativo
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
				}
			}
			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str
			estadoSeguimiento = helpers.ConsultarEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCualitativo", "err": "Error actualizando componente cualitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()
}

// GuardarCuantitativo ...
// @Title GuardarCuantitativo
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cuantitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCuantitativo() {
	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cuantitativo map[string]interface{}
	var informacion map[string]interface{}
	var estadoSeguimiento string
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
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
					panic(err)
				}
				detalle = map[string]interface{}{"estado": estado, "cuantitativo": cuantitativo, "informacion": informacion}
				dato[indiceActividad] = map[string]interface{}{"id": helpers.GuardarDetalleSeguimiento(detalle, false)}
			} else {
				identificador, segregado := dato[indiceActividad].(map[string]interface{})["id"]

				if segregado && identificador != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indiceActividad].(map[string]interface{})["id"].(string), &respuestaSeguimientoDetalle); err == nil {
						request.LimpiezaRespuestaRefactor(respuestaSeguimientoDetalle, &detalle)
						detalle = planeacion.ConvertirStringJson(detalle)
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

					observacion = helpers.ActividadConObservaciones(body)
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
			str := string(b)
			seguimiento["dato"] = str

			estadoSeguimiento = helpers.ConsultarEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()
}

// ReportarActividad ...
// @Title ReportarActividad
// @Description put Seguimiento by id
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_actividad/:index [put]
func (c *SeguimientoController) ReportarActividad() {
	indiceActividad := c.Ctx.Input.Param(":index")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var body map[string]interface{}
	var estado map[string]interface{}
	var respuestaSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+body["SeguimientoId"].(string), &respuesta); err == nil {
			aux := make(map[string]interface{}, 1)
			request.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)
			reportable, mensaje := helpers.ActividadReportable(seguimiento, indiceActividad)

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
						helpers.GuardarDetalleSeguimiento(detalle, true)
					}
				} else {
					dato[indiceActividad].(map[string]interface{})["estado"] = estado
					b, _ := json.Marshal(dato)
					str := string(b)
					seguimiento["dato"] = str

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				}
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimiento}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// ReportarSeguimiento ...
// @Title ReportarSeguimiento
// @Description put Seguimiento by id
// @Param	id			path 	string	true	"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_seguimiento/:id [put]
func (c *SeguimientoController) ReportarSeguimiento() {
	identificadorSeguimiento := c.Ctx.Input.Param(":id")
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+identificadorSeguimiento, &respuesta); err == nil {
		aux := make(map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux
		reportable, mensaje := helpers.SeguimientoReportable(seguimiento)

		if reportable {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:EAR", &respuestaEstado); err == nil {
				seguimiento["estado_seguimiento_id"] = respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
			}

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimiento}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// RevisarActividad ...
// @Title RevisarActividad
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /revision_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RevisarActividad() {
	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
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
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:CO", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				} else {
					panic(err)
				}
			} else {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AAV", &respuestaEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				} else {
					panic(err)
				}
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
					helpers.GuardarDetalleSeguimiento(detalle, true)
				}
			} else {
				dato[indiceActividad].(map[string]interface{})["estado"] = estado
				b, _ := json.Marshal(dato)
				str := string(b)
				seguimiento["dato"] = str
			}

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			}
			data := respuesta["Data"].(map[string]interface{})
			data["Observación"] = comentario
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :id is empty
// @router /revision_seguimiento/:id [put]
func (c *SeguimientoController) RevisarSeguimiento() {
	seguimientoIdentificador := c.Ctx.Input.Param(":id")
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
		avalado, observacion, mensaje := helpers.SeguimientoAvalable(seguimiento)

		if avalado || observacion {
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
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			}

			data := respuesta["Data"].(map[string]interface{})
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// RetornarActividad ...
// @Title RetornarActividad
// @Description Retorna la actividad de Avalado a en Revision
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /retornar_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RetornarActividad() {
	planIdentificador := c.Ctx.Input.Param(":plan_id")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planIdentificador+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
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
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				}

				if detalle["estado"].(map[string]interface{})["id"] == "63793207242b813898e9856b" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:OAPC", &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					} else {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					}
					seguimiento["estado_seguimiento_id"] = estado["id"]

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &respuestaEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     respuestaEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					} else {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					}
					detalle["estado"] = estado
					helpers.GuardarDetalleSeguimiento(detalle, true)

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					} else {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					}
					data := respuesta["Data"].(map[string]interface{})
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
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
					}
					dato[indiceActividad].(map[string]interface{})["estado"] = estado

					b, _ := json.Marshal(dato)
					str := string(b)
					seguimiento["dato"] = str

					if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					}
					data := respuesta["Data"].(map[string]interface{})
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}
				}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// MigrarInformacion ...
// @Title MigrarInformacion
// @Description post Segrar la informacion de los seguimientos
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /migrar_seguimiento/:plan_id/:trimestre [post]
func (c *SeguimientoController) MigrarInformacion() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	planIdentificador := c.Ctx.Input.Param(":plan_id")
	trimestre := c.Ctx.Input.Param(":trimestre")
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
			str := string(b)
			seguimiento["dato"] = str

			if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Actividades migradas:": respuestaMigrado, "Actividades no migradas: ": respuestaNoMigrado}}
	c.ServeJSON()
}
