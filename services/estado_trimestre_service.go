package services

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func EstadoTrimestres(planId string) (interface{}, error) {
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var auxPlanes []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId, &resSeguimiento); err == nil {
		request.LimpiezaRespuestaRefactor(resSeguimiento, &planes)

		for _, plan := range planes {
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

									auxPlanes = append(auxPlanes, plan)
								}
							}
						}
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
		if auxPlanes != nil {
			return auxPlanes, nil
		} else {
			return []int{}, nil
		}
	} else {
		logs.Error("Error -->", err)
		return nil, errors.New(err.Error())
	}
}

func ConsultarEstadoTrimestre(planIdentificador string, trimestre string) (interface{}, error) {
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
