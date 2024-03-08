package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestConsultarPeriodos(t *testing.T) {
	if response, err := http.Get("http://localhost:9013/v1/seguimiento/consultar_periodos/25"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPeriodos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPeriodos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPeriodos:", err.Error())
		t.Fail()
	}
}

func TestConsultarActividadesGenerales(t *testing.T) {
	if response, err := http.Get("http://localhost:9013/v1/seguimiento/consultar_actividades/61f60e4525e40c6f5d084185"); err == nil {
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

// SE NECESITAN DATOS PARA PODER VALIDAR EL CASO
func TestConsultarSeguimiento(t *testing.T) {
	if response, err := http.Get("http://localhost:9013/v1/seguimiento/consultar_seguimiento/61f08edc25e40c91b0083e4f/1/635b1f995073f2675157dc7f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarSeguimiento Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarSeguimiento:", err.Error())
		t.Fail()
	}
}

func TestConsultarIndicadores(t *testing.T) {
	if response, err := http.Get("http://localhost:9013/v1/seguimiento/consultar_indicadores/6201d43f25e40c205608b459"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarIndicadores Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarIndicadores Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarIndicadores:", err.Error())
		t.Fail()
	}
}

// SE NECESITAN DATOS PARA PODER VALIDAR EL CASO
func TestConsultarEstadoTrimestre(t *testing.T) {
	if response, err := http.Get("http://localhost:9013/v1/seguimiento/consultar_estado_trimestre/628ce817ebe1e6512a74b32e/T4"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarEstadoTrimestre Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarEstadoTrimestre Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarEstadoTrimestre:", err.Error())
		t.Fail()
	}
}

func TestCrearReportes(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:9013/v1/seguimiento/crear_reportes/61f08edc25e40c91b0083e4f/61f236f525e40c582a0840d0", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestClonarFormato Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestClonarFormato Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestClonarFormato:", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestConsultarAvanceIndicador(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:9013/v1/seguimiento/consultar_avance", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarAvanceIndicador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarAvanceIndicador Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarAvanceIndicador:", err.Error())
		t.Fail()
	}
}

func TestMigrarInformacion(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:9013/v1/seguimiento/migrar_seguimiento/prueba/635b1f795073f2675157dc7d", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestMigrarInformacion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestMigrarInformacion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestMigrarInformacion:", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestHabilitarReportes(t *testing.T) {
	body := []byte(`{
		"_id": "635b1f995073f2675157dc7f",
		"fecha_inicio": "2024-01-17T00:00:00.000Z",
		"fecha_fin": "2024-02-23T23:59:59.000Z",
		"periodo_id": "314",
		"activo": true,
		"fecha_creacion": "2022-10-28T00:17:29.116Z",
		"fecha_modificacion": "2024-03-04T14:42:36.388Z",
		"__v": 0,
		"tipo_seguimiento_id": "61f236f525e40c582a0840d0",
		"unidades_interes": "[{\"Id\":8,\"Nombre\":\"VICERRECTORIA ACADEMICA\"}]",
		"planes_interes": "[{\"_id\":\"628ce817ebe1e6512a74b32e\",\"nombre\":\"prueba nueva\"}]"
	  }`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/habilitar_reportes", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestHabilitarReportes Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestHabilitarReportes Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestGuardarSeguimiento(t *testing.T) {
	body := []byte(`{
		"_id": "64540e741287e015ea2f6d31",
		"cualitativo": "{\"dificultades\":\"No se apresentaron dificultades durante el tirmestre en el desarrollo de la actividad general. \",\"observaciones\":\"Obs cuali\",\"productos\":\"• Plan de Acción 2023 consolidado\",\"reporte\":\"Con el fin de consolidar el Plan de Acción 2023, la Oficina remitió durante la primera semana de enero los oficios respectivos informando la asignación presupuestal a cada Unidad y ordenador de gasto, a partir de las cuales se le solicitó a cada unidad ajustar la versión de su Plan de Acción. Ampliación de la descripción. \\nAtendiendo a su rol de acompañamiento, la Oficina realizó 34 sesiones, entre el 15 y 18 de marzo, orientadas a atender las inquietudes y observaciones de cada dependencia. Fruto de este trabajo, el 31 de marzo, se publicó el Plan de Acción de la Universidad en el portal web de la OAPC.  \"}",
		"cuantitativo": "{\"indicadores\":[{\"denominador\":\"Denominador fijo\",\"detalleReporte\":\"Se publicó el Plan de Acción Anual. \",\"formula\":\" Σ Planes e informes elaborados y publicados oportunamente de acuerdo con el cronograma definido por la OAPC\",\"meta\":\"5\",\"nombre\":\"Documentos (planes e informes) elaborados oportunamente\",\"observaciones\":\"ob cuant\",\"reporteDenominador\":\"1\",\"reporteNumerador\":\"1\",\"tendencia\":\"Creciente\",\"unidad\":\"Unidad\"},{\"denominador\":\"Denominador variable\",\"detalleReporte\":\"Denominador resulta de 34 sesiones promovidas por la OAPC, 0 sesiones solicitadas por las Unidades. \",\"formula\":\"(Sesiones de acompañamiento y asesoría realizadas/Solicitudes de acompañamiento de las Unidades Académicas y Administrativas + Sesiones programadas )*100\",\"meta\":100,\"nombre\":\"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción \",\"observaciones\":\"ob cuant 2\",\"reporteDenominador\":\"35\",\"reporteNumerador\":\"35\",\"tendencia\":\"Creciente\",\"unidad\":\"Porcentaje\"}],\"resultados\":[{\"acumuladoDenominador\":1,\"acumuladoNumerador\":1,\"avanceAcumulado\":0.2,\"brechaExistente\":4,\"divisionCero\":false,\"indicador\":1,\"indicadorAcumulado\":1,\"meta\":5,\"nombre\":\"Documentos (planes e informes) elaborados oportunamente\",\"unidad\":\"Unidad\"},{\"acumuladoDenominador\":35,\"acumuladoNumerador\":35,\"avanceAcumulado\":1,\"brechaExistente\":0,\"divisionCero\":false,\"indicador\":1,\"indicadorAcumulado\":1,\"meta\":100,\"nombre\":\"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción \",\"unidad\":\"Porcentaje\"}]}",
		"estado": "{\"id\":\"6361e98486e163cc4f1ed251\",\"nombre\":\"Actividad en reporte\"}",
		"evidencia": "[{\"Activo\":true,\"Enlace\":\"3c9711de-2602-4547-a567-2de493f8fa5e\",\"Id\":148487,\"Observacion\":\"\",\"TipoDocumento\":{\"codigoAbreviacion\":\"DSPA\",\"id\":60},\"nombre\":\"Gestión semanal Oficina Asesora de Planeación y Control - 6 al 10 de marzo.pdf\"},{\"Activo\":true,\"Enlace\":\"c31fd4c5-c12d-4b6b-9955-408bd3254562\",\"Id\":148488,\"Observacion\":\"\",\"TipoDocumento\":{\"codigoAbreviacion\":\"DSPA\",\"id\":60},\"nombre\":\"Gestión semanal Oficina Asesora de Planeación y Control - 13 al 17 de marzo.pdf\"}]",
		"informacion": "{\"descripcion\":\"Coordinar y asesorar metodológica y técnicamente el proceso de formulación, seguimiento y evaluación del Plan de Acción 2023, de acuerdo con lo establecido en el Sistema de Planeación Institucional. \",\"index\":\"1\",\"nombre\":\"Plan de acción 2023 Prod Seguimiento\",\"periodo\":\"Toda la vigencia\",\"ponderacion\":60,\"producto\":\"• Plan Operativo General 2023\\n• Plan de Acción 2023\\n• Informes de seguimiento al Plan de Acción 2023\",\"tarea\":\"• Comunicar el presupuesto aprobado a cada Unidad Académica y Administrativa, con el fin de que revisen y ajusten, en los casos que haya lugar, su Plan de Acción.\\n• Asesorar y acompañar técnicamente a las Unidades Académicas y Administrativas en el ajuste a su Plan de Acción. \\n• Consolidar el Plan Operativo General y el Plan de Acción 2023,a partir del Plan de Acción ajustado de cada Unidad.\\n• Publicar en el Portal Web de la Oficina Asesora de Planeación y Control el Plan de Acción y el Plan Operativo General 2023.\\n• Definir y socializar la metodología y herramienta para el seguimiento al Plan de Acción con las Unidades Académicas y Administrativas. \\n• Acompañar técnica y metodológicamente a las dependencias en el ejercicio de seguimiento a sus planes.\\n• Revisar y analizar los reportes de seguimiento de las Unidades y establecer observaciones y recomendaciones. \\n• Realizar la consolidación del seguimiento trimestral a los Planes de Acción y analizar los resultados asociados a cada trimestre. \\n• Generar el informe de seguimiento trimestral al Plan de Acción Anual.\\n• Realizar la publicación respectiva del documento generado.\\n• Diseñar espacios de socialización dirigidos a los jefes de las dependencias y los enlaces. \\n• Desplegar estrategias de socialización y divulgación de los resultados del ejercicio de seguimiento.\",\"trimestre\":\"T1\",\"unidad\":\"FACULTAD DE CIENCIAS MATEMATICAS Y NATURALES\"}",
		"fecha_creacion": "2023-05-04T19:58:44.598Z",
		"fecha_modificacion": "2023-05-04T19:58:44.598Z",
		"activo": true,
		"__v": 0
	  }`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/guardar_seguimiento/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarSeguimiento Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestGuardarDocumentos(t *testing.T) {
	body := []byte(`{   
		"file": "si",
		"documento": [{
        "IdTipoDocumento": 66,
          "nombre": "PRUEBA",
          "metadatos": {
            "dato_a": "Soportes planeacion"
          },
          "descripcion": "Documento de soporte para proyectos de plan de acción de inversión",
          "file": "DATA"
    	}]
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/guardar_documentos/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarDocumentos Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarDocumentos Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestGuardarCualitativo(t *testing.T) {
	body := []byte(`{   
		"cualitativo": {"dificultades":"No se apresentaron dificultades durante el tirmestre en el desarrollo de la actividad general. ","observaciones":"Obs cuali","productos":"• Plan de Acción 2023 consolidado","reporte":"Con el fin de consolidar el Plan de Acción 2023, la Oficina remitió durante la primera semana de enero los oficios respectivos informando la asignación presupuestal a cada Unidad y ordenador de gasto, a partir de las cuales se le solicitó a cada unidad ajustar la versión de su Plan de Acción. Ampliación de la descripción. \nAtendiendo a su rol de acompañamiento, la Oficina realizó 34 sesiones, entre el 15 y 18 de marzo, orientadas a atender las inquietudes y observaciones de cada dependencia. Fruto de este trabajo, el 31 de marzo, se publicó el Plan de Acción de la Universidad en el portal web de la OAPC.  "},
		"cuantitativo": {"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"Se publicó el Plan de Acción Anual. ","formula":" Σ Planes e informes elaborados y publicados oportunamente de acuerdo con el cronograma definido por la OAPC","meta":"5","nombre":"Documentos (planes e informes) elaborados oportunamente","observaciones":"ob cuant","reporteDenominador":"1","reporteNumerador":"1","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"Denominador variable","detalleReporte":"Denominador resulta de 34 sesiones promovidas por la OAPC, 0 sesiones solicitadas por las Unidades. ","formula":"(Sesiones de acompañamiento y asesoría realizadas/Solicitudes de acompañamiento de las Unidades Académicas y Administrativas + Sesiones programadas )*100","meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","observaciones":"ob cuant 2","reporteDenominador":"35","reporteNumerador":"35","tendencia":"Creciente","unidad":"Porcentaje"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":1,"avanceAcumulado":0.2,"brechaExistente":4,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":5,"nombre":"Documentos (planes e informes) elaborados oportunamente","unidad":"Unidad"},{"acumuladoDenominador":35,"acumuladoNumerador":35,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","unidad":"Porcentaje"}]},
		"informacion": {"descripcion":"Coordinar y asesorar metodológica y técnicamente el proceso de formulación, seguimiento y evaluación del Plan de Acción 2023, de acuerdo con lo establecido en el Sistema de Planeación Institucional. ","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":60,"producto":"• Plan Operativo General 2023\n• Plan de Acción 2023\n• Informes de seguimiento al Plan de Acción 2023","tarea":"• Comunicar el presupuesto aprobado a cada Unidad Académica y Administrativa, con el fin de que revisen y ajusten, en los casos que haya lugar, su Plan de Acción.\n• Asesorar y acompañar técnicamente a las Unidades Académicas y Administrativas en el ajuste a su Plan de Acción. \n• Consolidar el Plan Operativo General y el Plan de Acción 2023,a partir del Plan de Acción ajustado de cada Unidad.\n• Publicar en el Portal Web de la Oficina Asesora de Planeación y Control el Plan de Acción y el Plan Operativo General 2023.\n• Definir y socializar la metodología y herramienta para el seguimiento al Plan de Acción con las Unidades Académicas y Administrativas. \n• Acompañar técnica y metodológicamente a las dependencias en el ejercicio de seguimiento a sus planes.\n• Revisar y analizar los reportes de seguimiento de las Unidades y establecer observaciones y recomendaciones. \n• Realizar la consolidación del seguimiento trimestral a los Planes de Acción y analizar los resultados asociados a cada trimestre. \n• Generar el informe de seguimiento trimestral al Plan de Acción Anual.\n• Realizar la publicación respectiva del documento generado.\n• Diseñar espacios de socialización dirigidos a los jefes de las dependencias y los enlaces. \n• Desplegar estrategias de socialización y divulgación de los resultados del ejercicio de seguimiento.","trimestre":"T1","unidad":"FACULTAD DE CIENCIAS MATEMATICAS Y NATURALES"},
		"dependencia": "PRUEBA"
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/guardar_cualitativo/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarCualitativo Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarCualitativo Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestGuardarCuantitativo(t *testing.T) {
	body := []byte(`{   
		"cualitativo": {"dificultades":"No se apresentaron dificultades durante el tirmestre en el desarrollo de la actividad general. ","observaciones":"Obs cuali","productos":"• Plan de Acción 2023 consolidado","reporte":"Con el fin de consolidar el Plan de Acción 2023, la Oficina remitió durante la primera semana de enero los oficios respectivos informando la asignación presupuestal a cada Unidad y ordenador de gasto, a partir de las cuales se le solicitó a cada unidad ajustar la versión de su Plan de Acción. Ampliación de la descripción. \nAtendiendo a su rol de acompañamiento, la Oficina realizó 34 sesiones, entre el 15 y 18 de marzo, orientadas a atender las inquietudes y observaciones de cada dependencia. Fruto de este trabajo, el 31 de marzo, se publicó el Plan de Acción de la Universidad en el portal web de la OAPC.  "},
		"cuantitativo": {"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"Se publicó el Plan de Acción Anual. ","formula":" Σ Planes e informes elaborados y publicados oportunamente de acuerdo con el cronograma definido por la OAPC","meta":"5","nombre":"Documentos (planes e informes) elaborados oportunamente","observaciones":"ob cuant","reporteDenominador":"1","reporteNumerador":"1","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"Denominador variable","detalleReporte":"Denominador resulta de 34 sesiones promovidas por la OAPC, 0 sesiones solicitadas por las Unidades. ","formula":"(Sesiones de acompañamiento y asesoría realizadas/Solicitudes de acompañamiento de las Unidades Académicas y Administrativas + Sesiones programadas )*100","meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","observaciones":"ob cuant 2","reporteDenominador":"35","reporteNumerador":"35","tendencia":"Creciente","unidad":"Porcentaje"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":1,"avanceAcumulado":0.2,"brechaExistente":4,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":5,"nombre":"Documentos (planes e informes) elaborados oportunamente","unidad":"Unidad"},{"acumuladoDenominador":35,"acumuladoNumerador":35,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":100,"nombre":"Acompañamiento y asesoría en la formulación y seguimiento al Plan de Acción ","unidad":"Porcentaje"}]},
		"informacion": {"descripcion":"Coordinar y asesorar metodológica y técnicamente el proceso de formulación, seguimiento y evaluación del Plan de Acción 2023, de acuerdo con lo establecido en el Sistema de Planeación Institucional. ","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":60,"producto":"• Plan Operativo General 2023\n• Plan de Acción 2023\n• Informes de seguimiento al Plan de Acción 2023","tarea":"• Comunicar el presupuesto aprobado a cada Unidad Académica y Administrativa, con el fin de que revisen y ajusten, en los casos que haya lugar, su Plan de Acción.\n• Asesorar y acompañar técnicamente a las Unidades Académicas y Administrativas en el ajuste a su Plan de Acción. \n• Consolidar el Plan Operativo General y el Plan de Acción 2023,a partir del Plan de Acción ajustado de cada Unidad.\n• Publicar en el Portal Web de la Oficina Asesora de Planeación y Control el Plan de Acción y el Plan Operativo General 2023.\n• Definir y socializar la metodología y herramienta para el seguimiento al Plan de Acción con las Unidades Académicas y Administrativas. \n• Acompañar técnica y metodológicamente a las dependencias en el ejercicio de seguimiento a sus planes.\n• Revisar y analizar los reportes de seguimiento de las Unidades y establecer observaciones y recomendaciones. \n• Realizar la consolidación del seguimiento trimestral a los Planes de Acción y analizar los resultados asociados a cada trimestre. \n• Generar el informe de seguimiento trimestral al Plan de Acción Anual.\n• Realizar la publicación respectiva del documento generado.\n• Diseñar espacios de socialización dirigidos a los jefes de las dependencias y los enlaces. \n• Desplegar estrategias de socialización y divulgación de los resultados del ejercicio de seguimiento.","trimestre":"T1","unidad":"FACULTAD DE CIENCIAS MATEMATICAS Y NATURALES"},
		"dependencia": "PRUEBA"
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/guardar_cuantitativo/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarCuantitativo Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarCuantitativo Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// NO HAY DATA QUE CUMPLA CON LOS PARAMETROS DEL CONTROLADOR
func TestReportarActividad(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/reportar_actividad/:index", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestReportarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestReportarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// NO HAY DATA QUE CUMPLA CON LOS PARAMETROS DEL CONTROLADOR
func TestReportarSeguimiento(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/reportar_seguimiento/:id", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestReportarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestReportarSeguimiento Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
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

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/revision_actividad/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
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

func TestRevisarSeguimiento(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/revision_seguimiento/639a42e954a3d2399c3bb6ff", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestRevisarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestRevisarSeguimiento Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// SE NECESITAN EL JSON
func TestRetornarActividad(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:9013/v1/seguimiento/retornar_actividad/:plan_id/:index/:trimestre", bytes.NewBuffer(body)); err == nil {
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
