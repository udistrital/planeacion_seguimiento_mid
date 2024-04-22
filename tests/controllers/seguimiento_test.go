package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

// SE NECESITAN DATOS PARA PODER VALIDAR EL CASO
func TestConsultarSeguimiento(t *testing.T) {
	if response, err := http.Get("http://localhost:8080/v1/seguimiento/61f08edc25e40c91b0083e4f/1/635b1f995073f2675157dc7f"); err == nil {
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

func TestMigrarInformacion(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:8080/v1/seguimiento/migracion/61f08edc25e40c91b0083e4f/635b1f995073f2675157dc7f", "application/json", bytes.NewBuffer(body)); err == nil {
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

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/reporte/", bytes.NewBuffer(body)); err == nil {
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

func TestGuardarSeguimiento(t *testing.T) {
	body := []byte(`{"_id":"63b6ec61159830c1e0900956","cualitativo":{"dificultades":"t","observaciones":"","productos":"r","reporte":"k"},"cuantitativo":{"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"(25%*0.05+25%*0.3+25%*0.3+25%*0.05+25%*0.3)","formula":"∑ % avance en la tarea *ponderación de la tarea","meta":"80","nombre":"Avance en el Proyecto SGR Colciencias","observaciones":"","reporteDenominador":"1","reporteNumerador":"0.25","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ cursos impartidos durante la vigencia ","meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","observaciones":"","reporteDenominador":"1","reporteNumerador":"1","tendencia":"Creciente","unidad":"Unidad"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":0.25,"avanceAcumulado":0.3125,"brechaExistente":0.55,"divisionCero":false,"indicador":0.25,"indicadorAcumulado":0.25,"meta":80,"nombre":"Avance en el Proyecto SGR Colciencias","unidad":"Porcentaje"},{"acumuladoDenominador":1,"acumuladoNumerador":1,"avanceAcumulado":0.143,"brechaExistente":6,"divisionCero":false,"indicador":1,"indicadorAcumulado":1,"meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","unidad":"Unidad"}]},"estado":{"id":"63793207242b813898e9856b","nombre":"Actividad avalada"},"estadoSeguimiento":"En revisión OAPC","evidencia":[],"id":"63b6ec61159830c1e0900956","informacion":{"descripcion":"Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":25,"producto":"• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe","tarea":" Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)","trimestre":"T1","unidad":"FACULTAD DE INGENIERIA"}}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento/63b5fecb1598303a848fe7b8/1/635b1f205073f2675157dc7b", bytes.NewBuffer(body)); err == nil {
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

func TestRevisarSeguimiento(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento/639a42e954a3d2399c3bb6ff", bytes.NewBuffer(body)); err == nil {
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
