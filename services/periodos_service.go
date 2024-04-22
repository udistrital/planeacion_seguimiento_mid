package services

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func ConsultarTrimestres(vigencia string) ([]map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}
	var respuestaParametros map[string]interface{}

	if len(vigencia) == 0 {
		return nil, errors.New("error del servicio ConsultarTrimestres: Request containt incorrect params")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro?query=CodigoAbreviacion.in:T1|T2|T3|T4", &respuestaParametros); err == nil {
		var parametros []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaParametros, &parametros)

		for _, parametro := range parametros {
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:"+parametro["CodigoAbreviacion"].(string), &respuesta); err == nil {
				request.LimpiezaRespuestaRefactor(respuesta, &trimestre)
				trimestres = append(trimestres, trimestre...)
			} else {
				logs.Error("Error --> ", err)
				return nil, errors.New(err.Error())
			}
			trimestre = nil
		}
	} else {
		return nil, errors.New(err.Error())
	}
	return trimestres, nil
}
