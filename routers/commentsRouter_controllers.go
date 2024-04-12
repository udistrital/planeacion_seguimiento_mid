package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RetornarActividad",
            Router: "/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "ConsultarActividadesGenerales",
            Router: "/:seguimiento_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ActividadesController"],
        beego.ControllerComments{
            Method: "RevisarActividad",
            Router: "/revision/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:EstadoTrimestresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:EstadoTrimestresController"],
        beego.ControllerComments{
            Method: "EstadoTrimestres",
            Router: "/:plan_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:EstadoTrimestresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:EstadoTrimestresController"],
        beego.ControllerComments{
            Method: "ConsultarEstadoTrimestre",
            Router: "/:plan_id/:trimestre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:IndicadoresController"],
        beego.ControllerComments{
            Method: "ConsultarIndicadores",
            Router: "/:plan_id",
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

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"],
        beego.ControllerComments{
            Method: "HabilitarReportes",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"],
        beego.ControllerComments{
            Method: "CrearReportes",
            Router: "/:plan/:tipo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"],
        beego.ControllerComments{
            Method: "ReportarActividad",
            Router: "/actividad/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:ReporteController"],
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
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarSeguimiento",
            Router: "/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarSeguimiento",
            Router: "/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "MigrarInformacion",
            Router: "/migracion/:plan_id/:trimestre",
            AllowHTTPMethods: []string{"post"},
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

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"],
        beego.ControllerComments{
            Method: "GuardarCualitativo",
            Router: "/cualitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"],
        beego.ControllerComments{
            Method: "GuardarCuantitativo",
            Router: "/cuantitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoDetalleController"],
        beego.ControllerComments{
            Method: "GuardarDocumentos",
            Router: "/documento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
