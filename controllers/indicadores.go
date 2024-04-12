package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// IndicadoresController operations for Indicadores
type IndicadoresController struct {
	beego.Controller
}

// URLMapping ...
func (c *IndicadoresController) URLMapping() {
	c.Mapping("ConsultarIndicadores", c.ConsultarIndicadores)
	c.Mapping("ConsultarAvanceIndicador", c.ConsultarAvanceIndicador)
}

// ConsultarIndicadores ...
// @Title ConsultarIndicadores
// @Description get Seguimiento
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:plan_id [get]
func (c *IndicadoresController) ConsultarIndicadores() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan_identificador := c.Ctx.Input.Param(":plan_id")

	if resultado, err := services.ConsultarIndicadores(plan_identificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ConsultarAvanceIndicador ...
// @Title ConsultarAvanceIndicador
// @Description post Seguimiento by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Seguimiento
// @Failure 404
// @router /avance [post]
func (c *IndicadoresController) ConsultarAvanceIndicador() {
	defer errorhandler.HandlePanic(&c.Controller)

	if resultado, err := services.ConsultarAvanceIndicador(c.Ctx.Input.RequestBody); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}
