package helpers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func CambiarEstadoPlan(plan map[string]interface{}, idEstado string) (map[string]interface{}, error) {
	idPlan := plan["_id"].(string)
	plan["estado_plan_id"] = idEstado
	var res map[string]interface{}
	if err := request.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+idPlan, "PUT", &res, plan); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}
