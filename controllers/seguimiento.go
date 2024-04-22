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
	c.Mapping("ConsultarSeguimiento", c.ConsultarSeguimiento)
	c.Mapping("RevisarSeguimiento", c.RevisarSeguimiento)
	c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	c.Mapping("VerificarSeguimiento", c.VerificarSeguimiento)
	c.Mapping("MigrarInformacion", c.MigrarInformacion)
}

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /:planId/:index/:trimestre [put]
func (c *SeguimientoController) GuardarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	requestBody := c.Ctx.Input.RequestBody
	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.GuardarSeguimiento(requestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ConsultarSeguimiento ...
// @Title ConsultarSeguimiento
// @Description get Seguimiento
// @Param	planId 	path 	string	true		"The key for staticblock"
// @Param	index 	path 	string	true		"The key for staticblock"
// @Param	trimestre 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:planId/:index/:trimestre [get]
func (c *SeguimientoController) ConsultarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestreIdentificador := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.ConsultarSeguimiento(planIdentificador, indiceActividad, trimestreIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 404 :id is empty
// @router /:id [put]
func (c *SeguimientoController) RevisarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	seguimientoIdentificador := c.Ctx.Input.Param(":id")

	if resultado, err := services.RevisarSeguimiento(seguimientoIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// MigrarInformacion ...
// @Title MigrarInformacion
// @Description post Segrar la informacion de los seguimientos
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404
// @router /migracion/:planId/:trimestre [post]
func (c *SeguimientoController) MigrarInformacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.MigrarInformacion(planIdentificador, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
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
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}
