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
// @Param	seguimiento_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:seguimiento_id [get]
func (c *ActividadesController) ConsultarActividadesGenerales() {
	defer errorhandler.HandlePanic(&c.Controller)

	seguimiento_identificador := c.Ctx.Input.Param(":seguimiento_id")

	if resultado, err := services.ConsultarActividadesGenerales(seguimiento_identificador); err == nil {
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
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /revision/:plan_id/:index/:trimestre [put]
func (c *ActividadesController) RevisarActividad() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /:plan_id/:index/:trimestre [put]
func (c *ActividadesController) RetornarActividad() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
