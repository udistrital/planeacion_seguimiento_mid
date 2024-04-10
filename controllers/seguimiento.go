package controllers

import (
	"encoding/json"

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
	c.Mapping("HabilitarReportes", c.HabilitarReportes)
	c.Mapping("CrearReportes", c.CrearReportes)
	c.Mapping("ConsultarPeriodos", c.ConsultarPeriodos)
	c.Mapping("ConsultarActividadesGenerales", c.ConsultarActividadesGenerales)
	c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	c.Mapping("ConsultarSeguimiento", c.ConsultarSeguimiento)
	c.Mapping("ConsultarIndicadores", c.ConsultarIndicadores)
	c.Mapping("ConsultarAvanceIndicador", c.ConsultarAvanceIndicador)
	c.Mapping("ConsultarEstadoTrimestre", c.ConsultarEstadoTrimestre)
	c.Mapping("GuardarDocumentos", c.GuardarDocumentos)
	c.Mapping("GuardarCualitativo", c.GuardarCualitativo)
	c.Mapping("GuardarCuantitativo", c.GuardarCuantitativo)
	c.Mapping("ReportarActividad", c.ReportarActividad)
	c.Mapping("ReportarSeguimiento", c.ReportarSeguimiento)
	c.Mapping("RetornarActividad", c.RetornarActividad)
	c.Mapping("MigrarInformacion", c.MigrarInformacion)
	c.Mapping("EstadoTrimestres", c.EstadoTrimestres)
	c.Mapping("VerificarSeguimiento", c.VerificarSeguimiento)
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description put Seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403
// @router /habilitar_reportes [put]
func (c *SeguimientoController) HabilitarReportes() {
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
// @Success 200
// @Failure 403
// @router /crear_reportes/:plan/:tipo [post]
func (c *SeguimientoController) CrearReportes() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan_identificador := c.Ctx.Input.Param(":plan")
	tipo := c.Ctx.Input.Param(":tipo")

	resultado, err := services.CrearReportes(plan_identificador, tipo)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// ConsultarPeriodos ...
// @Title ConsultarPeriodos
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /consultar_periodos/:vigencia [get]
func (c *SeguimientoController) ConsultarPeriodos() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia := c.Ctx.Input.Param(":vigencia")

	if resultado, err := services.ConsultarTrimestres(vigencia); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// ConsultarActividadesGenerales ...
// @Title ConsultarActividadesGenerales
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_actividades/:seguimiento_id [get]
func (c *SeguimientoController) ConsultarActividadesGenerales() {
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

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /guardar_seguimiento/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	requestBody := c.Ctx.Input.RequestBody
	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_seguimiento/:plan_id/:index/:trimestre [get]
func (c *SeguimientoController) ConsultarSeguimiento() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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

// ConsultarIndicadores ...
// @Title ConsultarIndicadores
// @Description get Seguimiento
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /consultar_indicadores/:plan_id [get]
func (c *SeguimientoController) ConsultarIndicadores() {
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
// @Failure 403 :plan_id is empty
// @router /consultar_avance [post]
func (c *SeguimientoController) ConsultarAvanceIndicador() {
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

// ConsultarEstadoTrimestre ...
// @Title ConsultarEstadoTrimestre
// @Description get Seguimiento del trimestre correspondiente
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404 not found resource
// @router /consultar_estado_trimestre/:plan_id/:trimestre [get]
func (c *SeguimientoController) ConsultarEstadoTrimestre() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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

// GuardarDocumentos ...
// @Title GuardarDocumentos
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_documentos/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarDocumentos() {
	defer errorhandler.HandlePanic(&c.Controller)

	requestBody := c.Ctx.Input.RequestBody
	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cualitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCualitativo() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cuantitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCuantitativo() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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

// ReportarActividad ...
// @Title ReportarActividad
// @Description put Seguimiento by id
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /reportar_actividad/:index [put]
func (c *SeguimientoController) ReportarActividad() {
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

// ReportarSeguimiento ...
// @Title ReportarSeguimiento
// @Description put Seguimiento by id
// @Param	id			path 	string	true	"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 404
// @router /reportar_seguimiento/:id [put]
func (c *SeguimientoController) ReportarSeguimiento() {
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

// RevisarActividad ...
// @Title RevisarActividad
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /revision_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RevisarActividad() {
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

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :id is empty
// @router /revision_seguimiento/:id [put]
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

// RetornarActividad ...
// @Title RetornarActividad
// @Description Retorna la actividad de Avalado a en Revision
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /retornar_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RetornarActividad() {
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

// MigrarInformacion ...
// @Title MigrarInformacion
// @Description post Segrar la informacion de los seguimientos
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /migrar_seguimiento/:plan_id/:trimestre [post]
func (c *SeguimientoController) MigrarInformacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	planIdentificador := c.Ctx.Input.Param(":plan_id")
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
// @router /verificar_seguimiento/:id [put]
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

// EstadoTrimestres ...
// @Title EstadoTrimestres
// @Description get Seguimiento de los trimestres correspondientes
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404 not found resource
// @router /estado_trimestres/:plan_id [get]
func (c *SeguimientoController) EstadoTrimestres() {
	defer errorhandler.HandlePanic(&c.Controller)

	planId := c.Ctx.Input.Param(":plan_id")

	if resultado, err := services.EstadoTrimestres(planId); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}
