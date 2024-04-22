package controllers

import (
	"net/http"
	"testing"
)

// SE NECESITAN DATOS PARA PODER VALIDAR EL CASO
func TestConsultarEstadoTrimestre(t *testing.T) {
	if response, err := http.Get("http://localhost:8080/v1/estado-trimestre/628ce817ebe1e6512a74b32e/T4"); err == nil {
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
