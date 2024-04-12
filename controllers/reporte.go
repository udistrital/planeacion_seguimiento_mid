package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// ReporteController operations for SeguimientoReportes
type ReporteController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteController) URLMapping() {
	c.Mapping("CrearReportes", c.CrearReportes)
	c.Mapping("HabilitarReportes", c.HabilitarReportes)
	c.Mapping("ReportarSeguimiento", c.ReportarSeguimiento)
	c.Mapping("ReportarActividad", c.ReportarActividad)
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description put Seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 404
// @router / [put]
func (c *ReporteController) HabilitarReportes() {
	defer errorhandler.HandlePanic(&c.Controller)

	var entrada map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada) // entrada := c.Ctx.Input.RequestBody

	if reportes, err := services.HabilitarReportes(entrada); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, reportes)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// CrearReportes ...
// @Title CrearReportes
// @Description Post Seguimiento
// @Param	plan 	path 	string	true		"The key for staticblock"
// @Param	tipo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:plan/:tipo [post]
func (c *ReporteController) CrearReportes() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan_identificador := c.Ctx.Input.Param(":plan")
	tipo := c.Ctx.Input.Param(":tipo")

	if resultado, err := services.CrearReportes(plan_identificador, tipo); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ReportarSeguimiento ...
// @Title ReportarSeguimiento
// @Description put Seguimiento by id
// @Param	id			path 	string	true	"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /seguimiento/:id [put]
func (c *ReporteController) ReportarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	identificadorSeguimiento := c.Ctx.Input.Param(":id")

	if resultado, err := services.ReportarSeguimiento(identificadorSeguimiento); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ReportarActividad ...
// @Title ReportarActividad
// @Description put Seguimiento by id
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /actividad/:index [put]
func (c *ReporteController) ReportarActividad() {
	defer errorhandler.HandlePanic(&c.Controller)

	indiceActividad := c.Ctx.Input.Param(":index")

	if resultado, err := services.ReportarActividad(c.Ctx.Input.RequestBody, indiceActividad); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
