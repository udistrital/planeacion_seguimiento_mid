package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestConsultarActividadesGenerales(t *testing.T) {
	if response, err := http.Get("http://localhost:8080/v1/actividades/61f60e4525e40c6f5d084185"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarActividadesGenerales Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarActividadesGenerales Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarActividadesGenerales:", err.Error())
		t.Fail()
	}
}

func TestRevisarActividad(t *testing.T) {
	body := []byte(`{   
		"cualitativo": {"dificultades":"No se apresentaron dificultades durante el tirmestre en el desarrollo de la actividad general. ","observaciones":"Obs cuali","productos":"• Plan de Acción 2023 consolidado","reporte":"Con el fin de consolidar el Plan de Acción 2023, la Oficina remitió durante la primera semana de enero los oficios respectivos informando la asignación presupuestal a cada Unidad y ordenador de gasto, a partir de las cuales se le solicitó a cada unidad ajustar la versión de su Plan de Acción. Ampliación de la descripción. \nAtendiendo a su rol de acompañamiento, la Oficina realizó 34 sesiones, entre el 15 y 18 de marzo, orientadas a atender las inquietudes y observaciones de cada dependencia. Fruto de este trabajo, el 31 de marzo, se publicó el Plan de Acción de la Universidad en el portal web de la OAPC.  "},
		"cuantitativo": {"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"Se publicó el Plan de Acción Anual. ","formula":" Σ Planes e informes elaborados y publicados oportunamente de acuerdo con el cronograma definido por la OAPC","meta":"5","nombre":"Documentos (planes e informes) elaborados oportunamente","observaciones":"ob cuant","reporteDenominador":"1","reporteNumerador":"1","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"Denominador variable","detalleReporte":"Denominador resulta de 34 sesiones promovidas por la OAPC, 0 sesiones solicitadas por las Unidades. ","formula":"(Sesiones de acompañamiento y asesoría realizadas/Solicitudes de acompañamiento de las Unidades Académicas y Administrativas + Sesiones programadas )*100","meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","observaciones":"ob cuant 2","reporteDenominador":"35","reporteNumerador":"35","tendencia":"Creciente","unidad":"Porcentaje"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":1,"avanceAcumulado":0.2,"brechaExistente":4,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":5,"nombre":"Documentos (planes e informes) elaborados oportunamente","unidad":"Unidad"},{"acumuladoDenominador":35,"acumuladoNumerador":35,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","unidad":"Porcentaje"}]},
		"informacion": {"descripcion":"Coordinar y asesorar metodológica y técnicamente el proceso de formulación, seguimiento y evaluación del Plan de Acción 2023, de acuerdo con lo establecido en el Sistema de Planeación Institucional. ","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":60,"producto":"• Plan Operativo General 2023\n• Plan de Acción 2023\n• Informes de seguimiento al Plan de Acción 2023","tarea":"• Comunicar el presupuesto aprobado a cada Unidad Académica y Administrativa, con el fin de que revisen y ajusten, en los casos que haya lugar, su Plan de Acción.\n• Asesorar y acompañar técnicamente a las Unidades Académicas y Administrativas en el ajuste a su Plan de Acción. \n• Consolidar el Plan Operativo General y el Plan de Acción 2023,a partir del Plan de Acción ajustado de cada Unidad.\n• Publicar en el Portal Web de la Oficina Asesora de Planeación y Control el Plan de Acción y el Plan Operativo General 2023.\n• Definir y socializar la metodología y herramienta para el seguimiento al Plan de Acción con las Unidades Académicas y Administrativas. \n• Acompañar técnica y metodológicamente a las dependencias en el ejercicio de seguimiento a sus planes.\n• Revisar y analizar los reportes de seguimiento de las Unidades y establecer observaciones y recomendaciones. \n• Realizar la consolidación del seguimiento trimestral a los Planes de Acción y analizar los resultados asociados a cada trimestre. \n• Generar el informe de seguimiento trimestral al Plan de Acción Anual.\n• Realizar la publicación respectiva del documento generado.\n• Diseñar espacios de socialización dirigidos a los jefes de las dependencias y los enlaces. \n• Desplegar estrategias de socialización y divulgación de los resultados del ejercicio de seguimiento.","trimestre":"T1","unidad":"FACULTAD DE CIENCIAS MATEMATICAS Y NATURALES"},
		"dependencia": "PRUEBA",
        "evidencia": [{"Activo":true,"Enlace":"3c9711de-2602-4547-a567-2de493f8fa5e","Id":148487,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"Gestión semanal Oficina Asesora de Planeación y Control - 6 al 10 de marzo.pdf"},{"Activo":true,"Enlace":"c31fd4c5-c12d-4b6b-9955-408bd3254562","Id":148488,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"Gestión semanal Oficina Asesora de Planeación y Control - 13 al 17 de marzo.pdf"}]
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/actividades/revision/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestRevisarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestRevisarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestRetornarActividad(t *testing.T) {
	body := []byte(`{"_id":"63b6ec611598306b5c90095b","cualitativo":{"dificultades":"grdfvbfdesafghndawSHGFRDFRGTHYJKULIHGJUYTRFEHJVNGBFEW3R4T5Y67UTJYHFGRCUTYRFESRGTHYJUGHGjhkjgfrdgh","observaciones":"","productos":"gbvnmbhgfvdsafghjm","reporte":"Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:"},"cuantitativo":{"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ % avance en la tarea *ponderación de la tarea","meta":"80","nombre":"Avance en el Proyecto SGR Colciencias","observaciones":"","reporteDenominador":"1","reporteNumerador":"0.3","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ cursos impartidos durante la vigencia ","meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","observaciones":"","reporteDenominador":"1","reporteNumerador":"2","tendencia":"Creciente","unidad":"Unidad"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":0.55,"avanceAcumulado":0.6875,"brechaExistente":0.25,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.55,"meta":80,"nombre":"Avance en el Proyecto SGR Colciencias","unidad":"Porcentaje"},{"acumuladoDenominador":1,"acumuladoNumerador":3,"avanceAcumulado":0.429,"brechaExistente":4,"divisionCero":false,"indicador":2,"indicadorAcumulado":3,"meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","unidad":"Unidad"}]},"estado":{"id":"63793207242b813898e9856b","nombre":"Actividad avalada"},"estadoSeguimiento":"En revisión OAPC","evidencia":[{"Activo":true,"Enlace":"c275df0d-b27a-4446-b434-69ee6b71d2c1","Id":148376,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"UNAL - 230210 - SOSTENIBILIDAD ECONÓMICA (CC).docx"},{"Activo":true,"Enlace":"927b1653-7c23-40f1-bd04-e35a98e9e72a","Id":148377,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"Acreditación de programas.xlsx"}],"id":"","informacion":{"descripcion":"Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":25,"producto":"• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe","tarea":" Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)","trimestre":"T2","unidad":"FACULTAD DE INGENIERIA"}}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/actividades/63b5fecb1598303a848fe7b8/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestRetornarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestRetornarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}
