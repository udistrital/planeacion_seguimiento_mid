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

	if resultado, err := services.SolicitudReformulacion(c.Ctx.Input.RequestBody); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
