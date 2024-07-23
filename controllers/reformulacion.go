package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// ReformulacionController operations for Reformulacion
type ReformulacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReformulacionController) URLMapping() {
	c.Mapping("SolicitudReformulacion", c.SolicitudReformulacion)
}

// SolicitudReformulacion ...
// @Title SolicitudReformulacion
// @Description Solicitud de reformulacion para un plan de acción
// @Param	body		body 	{}	true		"body for reformulacion content"
// @Success 201
// @Failure 404
// @router / [post]
func (c *ReformulacionController) SolicitudReformulacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	if resultado, err := services.SolicitarReformulacion(c.Ctx.Input.RequestBody); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 201, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ValidacionReformulacion ...
// @Title ValidacionReformulacion
// @Description Validación de la habilidad para realizar una reformulacion para un plan de acción
// @Param	body		body 	{}	true		"body for reformulacion content"
// @Success 200
// @Failure 404
// @router /validar/:plan_id [get]
func (c *ReformulacionController) ValidacionReformulacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")

	if resultado, err := services.ValidacionReformulacion(planIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 400, nil, err.Error())
	}
	c.ServeJSON()
}

// AprobarReformulacion ...
// @Title AprobarReformulacion
// @Description Aprobación de una reformulacion para un plan de acción
// @Param	body		body 	{}	true		"body for reformulacion content"
// @Success 200
// @Failure 404
// @router /validar/:reformulacion_id [get]
func (c *ReformulacionController) AprobarReformulacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":reformulacion_id")

	if resultado, err := services.AprobarReformulacion(planIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 500, nil, err.Error())
	}
	c.ServeJSON()
}
