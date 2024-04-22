package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

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
		"evidencia": []
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento-detalle/documento/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
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
		"_id": "",
		"informacion": {
			"descripcion": "Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.",
			"index": "1",
			"nombre": "Plan de acción 2023 Prod Seguimiento",
			"periodo": "Toda la vigencia",
			"ponderacion": 25,
			"producto": "• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe",
			"tarea": " Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)",
			"trimestre": "T2",
			"unidad": "FACULTAD DE INGENIERIA"
		},
		"evidencias": [
			{
				"Activo": true,
				"Enlace": "c275df0d-b27a-4446-b434-69ee6b71d2c1",
				"Id": 148376,
				"Observacion": "",
				"TipoDocumento": {
					"codigoAbreviacion": "DSPA",
					"id": 60
				},
				"nombre": "UNAL - 230210 - SOSTENIBILIDAD ECONÓMICA (CC).docx"
			},
			{
				"Activo": true,
				"Enlace": "927b1653-7c23-40f1-bd04-e35a98e9e72a",
				"Id": 148377,
				"Observacion": "",
				"TipoDocumento": {
					"codigoAbreviacion": "DSPA",
					"id": 60
				},
				"nombre": "Acreditación de programas.xlsx"
			}
		],
		"cualitativo": {
			"dificultades": "grdfvbfdesafghndawSHGFRDFRGTHYJKULIHGJUYTRFEHJVNGBFEW3R4T5Y67UTJYHFGRCUTYRFESRGTHYJUGHGjhkjgfrdgh",
			"observaciones": "",
			"productos": "gbvnmbhgfvdsafghjm",
			"reporte": "Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:"
		},
		"cuantitativo": {
			"indicadores": [
				{
					"denominador": "Denominador fijo",
					"detalleReporte": "",
					"formula": "∑ % avance en la tarea *ponderación de la tarea",
					"meta": "80",
					"nombre": "Avance en el Proyecto SGR Colciencias",
					"observaciones": "",
					"reporteDenominador": "1",
					"reporteNumerador": "0.3",
					"tendencia": "Creciente",
					"unidad": "Porcentaje"
				},
				{
					"denominador": "Denominador fijo",
					"detalleReporte": "",
					"formula": "∑ cursos impartidos durante la vigencia ",
					"meta": 7,
					"nombre": "Cursos Cultiva articulados a la Red Acacia",
					"observaciones": "",
					"reporteDenominador": "1",
					"reporteNumerador": "2",
					"tendencia": "Creciente",
					"unidad": "Unidad"
				}
			],
			"resultados": [
				{
					"acumuladoDenominador": 1,
					"acumuladoNumerador": 0.55,
					"avanceAcumulado": 0.6875,
					"brechaExistente": 0.25,
					"divisionCero": false,
					"indicador": 0.3,
					"indicadorAcumulado": 0.55,
					"meta": 80,
					"nombre": "Avance en el Proyecto SGR Colciencias",
					"unidad": "Porcentaje"
				},
				{
					"acumuladoDenominador": 1,
					"acumuladoNumerador": 3,
					"avanceAcumulado": 0.429,
					"brechaExistente": 4,
					"divisionCero": false,
					"indicador": 2,
					"indicadorAcumulado": 3,
					"meta": 7,
					"nombre": "Cursos Cultiva articulados a la Red Acacia",
					"unidad": "Unidad"
				}
			]
		},
		"dependencia": false
	}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento-detalle/cualitativo/63b5fecb1598303a848fe7b8/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
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
	body := []byte(`{"_id":"","informacion":{"descripcion":"Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":25,"producto":"• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe","tarea":" Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)","trimestre":"T2","unidad":"FACULTAD DE INGENIERIA"},"evidencias":[{"Activo":true,"Enlace":"c275df0d-b27a-4446-b434-69ee6b71d2c1","Id":148376,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"UNAL - 230210 - SOSTENIBILIDAD ECONÓMICA (CC).docx"},{"Activo":true,"Enlace":"927b1653-7c23-40f1-bd04-e35a98e9e72a","Id":148377,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"Acreditación de programas.xlsx"}],"cualitativo":{"dificultades":"grdfvbfdesafghndawSHGFRDFRGTHYJKULIHGJUYTRFEHJVNGBFEW3R4T5Y67UTJYHFGRCUTYRFESRGTHYJUGHGjhkjgfrdgh","observaciones":"","productos":"gbvnmbhgfvdsafghjm","reporte":"Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:"},"cuantitativo":{"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ % avance en la tarea *ponderación de la tarea","meta":"80","nombre":"Avance en el Proyecto SGR Colciencias","observaciones":"","reporteDenominador":"1","reporteNumerador":"0.3","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ cursos impartidos durante la vigencia ","meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","observaciones":"","reporteDenominador":"1","reporteNumerador":"2","tendencia":"Creciente","unidad":"Unidad"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":0.55,"avanceAcumulado":0.6875,"brechaExistente":0.25,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.55,"meta":80,"nombre":"Avance en el Proyecto SGR Colciencias","unidad":"Porcentaje"},{"acumuladoDenominador":1,"acumuladoNumerador":3,"avanceAcumulado":0.429,"brechaExistente":4,"divisionCero":false,"indicador":2,"indicadorAcumulado":3,"meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","unidad":"Unidad"}]},"dependencia":false}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento-detalle/cuantitativo/63b5fecb1598303a848fe7b8/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
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
