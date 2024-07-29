package services

import (
	"errors"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"golang.org/x/sync/errgroup"
)

func ConsultarTrimestres(vigencia string) ([]map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var respuestaParametros map[string]interface{}
	var mutex sync.Mutex
	wge := new(errgroup.Group)

	if len(vigencia) == 0 {
		return nil, errors.New("error del servicio ConsultarTrimestres: Request containt incorrect params")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro?query=CodigoAbreviacion.in:T1|T2|T3|T4", &respuestaParametros); err != nil {
		return nil, errors.New(err.Error())
	}

	var parametros []map[string]interface{}
	request.LimpiezaRespuestaRefactor(respuestaParametros, &parametros)

	// Mapa para almacenar resultados de los trimestres
	trimestreMap := make(map[string][]map[string]interface{})

	for _, parametro := range parametros {
		parametro := parametro
		wge.Go(func() error {
			var trimestre []map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:"+parametro["CodigoAbreviacion"].(string), &respuesta); err != nil {
				logs.Error("Error --> ", err)
				return errors.New(err.Error())
			}

			request.LimpiezaRespuestaRefactor(respuesta, &trimestre)
			if len(trimestre) == 0 {
				return errors.New("error del servicio ConsultarTrimestres: No se encontraron trimestres para la vigencia " + vigencia)
			}

			mutex.Lock()
			trimestreMap[parametro["CodigoAbreviacion"].(string)] = trimestre
			mutex.Unlock()

			return nil
		})
	}

	if err := wge.Wait(); err != nil {
		return nil, errors.New(err.Error())
	}

	// Slice para almacenar los trimestres en orden
	var trimestres []map[string]interface{}

	for _, code := range []string{"T1", "T2", "T3", "T4"} {
		if trimes, exists := trimestreMap[code]; exists {
			trimestres = append(trimestres, trimes...)
		}
	}

	return trimestres, nil

	// var respuesta map[string]interface{}
	// var trimestre []map[string]interface{}
	// var trimestres []map[string]interface{}
	// var respuestaParametros map[string]interface{}
	// wge := new(errgroup.Group)
	// var mutex sync.Mutex
	//
	// if len(vigencia) == 0 {
	// 	return nil, errors.New("error del servicio ConsultarTrimestres: Request containt incorrect params")
	// }
	//
	// if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro?query=CodigoAbreviacion.in:T1|T2|T3|T4", &respuestaParametros); err == nil {
	// 	var parametros []map[string]interface{}
	// 	request.LimpiezaRespuestaRefactor(respuestaParametros, &parametros)
	//
	// 	for _, parametro := range parametros {
	// 		parametro := parametro
	// 		wge.Go(func() error {
	// 			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:"+parametro["CodigoAbreviacion"].(string), &respuesta); err == nil {
	// 				request.LimpiezaRespuestaRefactor(respuesta, &trimestre)
	// 				if len(trimestre[0]) == 0 {
	// 					return errors.New("error del servicio ConsultarTrimestres: No se encontraron trimestres para la vigencia " + vigencia)
	// 				}
	// 				mutex.Lock()
	// 				trimestres = append(trimestres, trimestre...)
	// 				mutex.Unlock()
	// 			} else {
	// 				logs.Error("Error --> ", err)
	// 				return errors.New(err.Error())
	// 			}
	// 			mutex.Lock()
	// 			trimestre = nil
	// 			mutex.Unlock()
	// 			return nil
	// 		})
	// 	}
	//
	// 	if err := wge.Wait(); err != nil {
	// 		return nil, errors.New(err.Error())
	// 	}
	//
	// } else {
	// 	return nil, errors.New(err.Error())
	// }
	//
	// return trimestres, nil
}
