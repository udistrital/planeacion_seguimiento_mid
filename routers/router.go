// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_seguimiento_mid/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/seguimiento",
			beego.NSInclude(
				&controllers.SeguimientoController{},
			),
		),
		beego.NSNamespace("/detalles",
			beego.NSInclude(
				&controllers.DetallesController{},
			),
		),
		beego.NSNamespace("/actividades",
			beego.NSInclude(
				&controllers.ActividadesController{},
			),
		),
		beego.NSNamespace("/reportes",
			beego.NSInclude(
				&controllers.ReportesController{},
			),
		),
		beego.NSNamespace("/indicadores",
			beego.NSInclude(
				&controllers.IndicadoresController{},
			),
		),
		beego.NSNamespace("/periodos",
			beego.NSInclude(
				&controllers.PeriodosController{},
			),
		),
		beego.NSNamespace("/reformulacion",
			beego.NSInclude(
				&controllers.ReformulacionController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
