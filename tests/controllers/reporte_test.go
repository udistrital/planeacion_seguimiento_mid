package controllers

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCrearReportes(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:8080/v1/reporte/61f08edc25e40c91b0083e4f/61f236f525e40c582a0840d0", "application/json", bytes.NewBuffer(body)); err == nil {
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

// NO HAY DATA QUE CUMPLA CON LOS PARAMETROS DEL CONTROLADOR
func TestReportarActividad(t *testing.T) {
	body := []byte(`{"SeguimientoId": "61f238db25e40ccb450840db"}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/reporte/actividad/1", bytes.NewBuffer(body)); err == nil {
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

func TestReportarSeguimiento(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/reporte/seguimiento/61f238db25e40ccb450840db", bytes.NewBuffer(body)); err == nil {
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
