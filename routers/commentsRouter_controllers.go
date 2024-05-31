package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "ConsultarActividadesGenerales",
            Router: "/:seguimientoId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RetornarActividad",
            Router: "/retornar/:planId/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RetornarActividadJefeDependencia",
            Router: "/retornar_jefe_dependencia/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RevisarActividad",
            Router: "/revision/:planId/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RevisarActividadJefeDependencia",
            Router: "/revision_jefe_dependencia/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"],
        beego.ControllerComments{
            Method: "GuardarCualitativo",
            Router: "/cualitativo/:planId/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"],
        beego.ControllerComments{
            Method: "GuardarCuantitativo",
            Router: "/cuantitativo/:planId/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:DetallesController"],
        beego.ControllerComments{
            Method: "GuardarDocumentos",
            Router: "/documento/:planId/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"],
        beego.ControllerComments{
            Method: "ConsultarIndicadores",
            Router: "/:planId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"],
        beego.ControllerComments{
            Method: "ConsultarAvanceIndicador",
            Router: "/avance",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:PeriodosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:PeriodosController"],
        beego.ControllerComments{
            Method: "ConsultarPeriodos",
            Router: "/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "CrearReportes",
            Router: "/:plan/:tipo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "ReportarActividad",
            Router: "/actividad/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "HabilitarReportes",
            Router: "/habilitar",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "ReportarSeguimiento",
            Router: "/seguimiento/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarSeguimiento",
            Router: "/:id/revision",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarSeguimiento",
            Router: "/:planId/:indiceActividad/:trimestreId",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarSeguimiento",
            Router: "/:planId/:indiceActividad/:trimestreId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarEstadoTrimestre",
            Router: "/:planId/:trimestre/estado",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "MigrarInformacion",
            Router: "/:planId/:trimestreId/migracion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "EstadoTrimestres",
            Router: "/:planId/estado",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "AvalarPlan",
            Router: "/avalar/:plan_id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarSeguimientoJefeDependencia",
            Router: "/revision_jefe_dependencia/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "VerificarSeguimiento",
            Router: "/verificacion/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
