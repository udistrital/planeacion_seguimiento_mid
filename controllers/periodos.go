package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// PeriodosController operations for Periodos
type PeriodosController struct {
	beego.Controller
}

// URLMapping ...
func (c *PeriodosController) URLMapping() {
	c.Mapping("ConsultarPeriodos", c.ConsultarPeriodos)
}

// ConsultarPeriodos ...
// @Title ConsultarPeriodos
// @Description get Seguimiento
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:vigencia [get]
func (c *PeriodosController) ConsultarPeriodos() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia := c.Ctx.Input.Param(":vigencia")

	if resultado, err := services.ConsultarTrimestres(vigencia); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}
