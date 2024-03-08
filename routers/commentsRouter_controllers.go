package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarActividadesGenerales",
            Router: "/consultar_actividades/:seguimiento_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarAvanceIndicador",
            Router: "/consultar_avance",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarEstadoTrimestre",
            Router: "/consultar_estado_trimestre/:plan_id/:trimestre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarIndicadores",
            Router: "/consultar_indicadores/:plan_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarPeriodos",
            Router: "/consultar_periodos/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ConsultarSeguimiento",
            Router: "/consultar_seguimiento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "CrearReportes",
            Router: "/crear_reportes/:plan/:tipo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarCualitativo",
            Router: "/guardar_cualitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarCuantitativo",
            Router: "/guardar_cuantitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarDocumentos",
            Router: "/guardar_documentos/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarSeguimiento",
            Router: "/guardar_seguimiento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "HabilitarReportes",
            Router: "/habilitar_reportes",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "MigrarInformacion",
            Router: "/migrar_seguimiento/:plan_id/:trimestre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ReportarActividad",
            Router: "/reportar_actividad/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ReportarSeguimiento",
            Router: "/reportar_seguimiento/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RetornarActividad",
            Router: "/retornar_actividad/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarActividad",
            Router: "/revision_actividad/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_seguimiento_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarSeguimiento",
            Router: "/revision_seguimiento/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
