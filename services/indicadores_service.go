package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func ConsultarIndicadores(plan_identificador string) (interface{}, error) {
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

					return indicadores, nil
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
	return nil, nil
}

func ConsultarAvanceIndicador(requestBody []byte) (interface{}, error) {
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
	json.Unmarshal(requestBody, &body)

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
							logs.Error("Error --> ", err)
						}
						periodoIdentificador = priodoId_rest - 1
					} else {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 32)
						if err != nil {
							logs.Error("Error --> ", err)
						}
						periodoIdentificador = priodoId_rest
					}
				} else {
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
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
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
				}

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_seguimiento_id"].(string), &respuestaName); err == nil {
					request.LimpiezaRespuestaRefactor(respuestaName, &parametro_periodo_name)
					paramIdlenName := parametro_periodo_name[0]

					paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
					nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
				} else {
					logs.Error("Error --> ", err)
					return nil, errors.New(err.Error())
				}
			}
		} else {
			logs.Error("Error --> ", err)
			return nil, errors.New(err.Error())
		}
		avancePeriodo := body["avancePeriodo"].(string)
		aPe, err := strconv.ParseFloat(avancePeriodo, 32)
		if err != nil {
			logs.Error("Error --> ", err)
		}

		aAc, err := strconv.ParseFloat(avanceAcumulado, 32)
		if err != nil {
			logs.Error("Error --> ", err)
		}
		totalAcumulado := fmt.Sprint(aPe + aAc)
		generalData := make(map[string]interface{})
		generalData["avancePeriodo"] = avancePeriodo
		generalData["periodIdString"] = periodoIdentificadorString
		generalData["avanceAcumulado"] = totalAcumulado
		generalData["avancePeriodoPrev"] = testavancePeriodo
		generalData["avanceAcumuladoPrev"] = avanceAcumulado
		generalData["nombrePeriodo"] = nombrePeriodo
		return generalData, nil
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New(err.Error())
	}
}
