package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// ActividadesController operations for ActividadesSeguimiento
type ActividadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActividadesController) URLMapping() {
	c.Mapping("ConsultarActividadesGenerales", c.ConsultarActividadesGenerales)
	c.Mapping("RetornarActividad", c.RetornarActividad)
	c.Mapping("RevisarActividad", c.RevisarActividad)
}

// ConsultarActividadesGenerales ...
// @Title ConsultarActividadesGenerales
// @Description get Seguimiento
// @Param	seguimientoId 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:seguimientoId [get]
func (c *ActividadesController) ConsultarActividadesGenerales() {
	defer errorhandler.HandlePanic(&c.Controller)

	seguimientoIdentificador := c.Ctx.Input.Param(":seguimientoId")

	if resultado, err := services.ConsultarActividadesGenerales(seguimientoIdentificador); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// RevisarActividad ...
// @Title RevisarActividad
// @Description put Seguimiento by id
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /revision/:planId/:index/:trimestre [put]
func (c *ActividadesController) RevisarActividad() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")
	if resultado, err := services.RevisarActividad(c.Ctx.Input.RequestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// RetornarActividad ...
// @Title RetornarActividad
// @Description Retorna la actividad de Avalado a en Revision
// @Param	planId		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /:planId/:index/:trimestre [put]
func (c *ActividadesController) RetornarActividad() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":planId")
	indiceActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")
	if resultado, err := services.RetornarActividad(c.Ctx.Input.RequestBody, planIdentificador, indiceActividad, trimestre); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
