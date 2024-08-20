package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// SeguimientoController operations for Seguimiento
type SeguimientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SeguimientoController) URLMapping() {
	c.Mapping("ObtenerSeguimientos", c.ObtenerSeguimientos)
	c.Mapping("ConsultarSeguimiento", c.ConsultarSeguimiento)
	c.Mapping("RevisarSeguimiento", c.RevisarSeguimiento)
	c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	c.Mapping("VerificarSeguimiento", c.VerificarSeguimiento)
	c.Mapping("MigrarInformacion", c.MigrarInformacion)
	c.Mapping("ConsultarEstadoTrimestre", c.ConsultarEstadoTrimestre)
	c.Mapping("EstadoTrimestres", c.EstadoTrimestres)
	c.Mapping("AvalarPlan", c.AvalarPlan)
	c.Mapping("RevisarSeguimientoJefeDependencia", c.RevisarSeguimientoJefeDependencia)
}

// ObtenerSeguimientos ...
// @Title ObtenerSeguimientos
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /:planId [get]
func (c *SeguimientoController) ObtenerSeguimientos() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")

	if resultado, err := services.ObtenerSeguimientos(planIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	indiceActividad		path 	string	true		"The key for staticblock"
// @Param	trimestreId	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /:planId/:indiceActividad/:trimestreId [put]
func (c *SeguimientoController) GuardarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	requestBody := c.Ctx.Input.RequestBody
	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":indiceActividad")
	trimestreId := c.Ctx.Input.Param(":trimestreId")

	if resultado, err := services.GuardarSeguimiento(requestBody, planIdentificador, indiceActividad, trimestreId); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ConsultarSeguimiento ...
// @Title ConsultarSeguimiento
// @Description get Seguimiento
// @Param	planId 				path 	string	true		"The key for staticblock"
// @Param	indiceActividad 	path 	string	true		"The key for staticblock"
// @Param	trimestreId 		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:planId/:indiceActividad/:trimestreId [get]
func (c *SeguimientoController) ConsultarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":indiceActividad")
	trimestreIdIdentificador := c.Ctx.Input.Param(":trimestreId")

	if resultado, err := services.ConsultarSeguimiento(planIdentificador, indiceActividad, trimestreIdIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 404 :id is empty
// @router /:id/revision [put]
func (c *SeguimientoController) RevisarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	seguimientoIdentificador := c.Ctx.Input.Param(":id")

	if resultado, err := services.RevisarSeguimiento(seguimientoIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// MigrarInformacion ...
// @Title MigrarInformacion
// @Description post Segrar la informacion de los seguimientos
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	trimestreId	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404
// @router /:planId/:trimestreId/migracion [post]
func (c *SeguimientoController) MigrarInformacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	trimestreId := c.Ctx.Input.Param(":trimestreId")

	if resultado, err := services.MigrarInformacion(planIdentificador, trimestreId); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// VerificarSeguimiento ...
// @Title VerificarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @Failure 404
// @router /verificacion/:id [put]
func (c *SeguimientoController) VerificarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	idSeguimiento := c.Ctx.Input.Param(":id")

	if resultado, err := services.VerificarSeguimiento(idSeguimiento); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// EstadoTrimestres ...
// @Title EstadoTrimestres
// @Description get Seguimiento de los trimestres correspondientes
// @Param	planId 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404 not found resource
// @router /:planId/estado [get]
func (c *SeguimientoController) EstadoTrimestres() {
	defer errorhandler.HandlePanic(&c.Controller)

	planId := c.Ctx.Input.Param(":planId")

	if resultado, err := services.EstadoTrimestres(planId); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
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
// @router /:planId/:trimestre/estado [get]
func (c *SeguimientoController) ConsultarEstadoTrimestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.ConsultarEstadoTrimestre(planIdentificador, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// AvalarPlan ...
// @Title AvalarPlan
// @Description Petición Post para avalar plan y crear reportes de seguimiento
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 400
// @router /avalar/:plan_id [post]
func (c *SeguimientoController) AvalarPlan() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan_id := c.Ctx.Input.Param(":plan_id")

	if resultado, err := services.AvalarPlan(plan_id); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :id is empty
// @router /revision_jefe_dependencia/:id [put]
func (c *SeguimientoController) RevisarSeguimientoJefeDependencia() {
	defer errorhandler.HandlePanic(&c.Controller)
	seguimiento_id := c.Ctx.Input.Param(":id")
	if resultado, err := services.RevisarSeguimientoJefeDependencia(seguimiento_id); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	c.ServeJSON()
}
