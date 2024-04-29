package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// DetallesController operations for SeguimientoDetalle
type DetallesController struct {
	beego.Controller
}

// URLMapping ...
func (c *DetallesController) URLMapping() {
	c.Mapping("GuardarDocumentos", c.GuardarDocumentos)
	c.Mapping("GuardarCualitativo", c.GuardarCualitativo)
	c.Mapping("GuardarCuantitativo", c.GuardarCuantitativo)
}

// GuardarDocumentos ...
// @Title GuardarDocumentos
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /documento/:planId/:index/:trimestre [put]
func (c *DetallesController) GuardarDocumentos() {
	defer errorhandler.HandlePanic(&c.Controller)

	requestBody := c.Ctx.Input.RequestBody
	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.GuardarDocumentos(requestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GuardarCualitativo ...
// @Title GuardarCualitativo
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /cualitativo/:planId/:index/:trimestre [put]
func (c *DetallesController) GuardarCualitativo() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.GuardarCualitativo(c.Ctx.Input.RequestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GuardarCuantitativo ...
// @Title GuardarCuantitativo
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /cuantitativo/:planId/:index/:trimestre [put]
func (c *DetallesController) GuardarCuantitativo() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if resultado, err := services.GuardarCuantitativo(c.Ctx.Input.RequestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
