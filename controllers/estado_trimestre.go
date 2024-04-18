package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EstadoTrimestresController operations for Trimestres
type EstadoTrimestresController struct {
	beego.Controller
}

// URLMapping ...
func (c *EstadoTrimestresController) URLMapping() {
	c.Mapping("ConsultarEstadoTrimestre", c.ConsultarEstadoTrimestre)
	c.Mapping("EstadoTrimestres", c.EstadoTrimestres)

}

// EstadoTrimestres ...
// @Title EstadoTrimestres
// @Description get Seguimiento de los trimestres correspondientes
// @Param	planId 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404 not found resource
// @router /:planId [get]
func (c *EstadoTrimestresController) EstadoTrimestres() {
	defer errorhandler.HandlePanic(&c.Controller)

	planId := c.Ctx.Input.Param(":planId")

	if resultado, err := services.EstadoTrimestres(planId); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ConsultarEstadoTrimestre ...
// @Title ConsultarEstadoTrimestre
// @Description get Seguimiento del trimestre correspondiente
// @Param	planId 	path 	string	true		"The key for staticblock"
// @Param	trimestre 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404 not found resource
// @router /:planId/:trimestre [get]
func (c *EstadoTrimestresController) ConsultarEstadoTrimestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.ConsultarEstadoTrimestre(planIdentificador, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}
