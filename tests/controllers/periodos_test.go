package controllers

import (
	"net/http"
	"testing"
)

func TestConsultarPeriodos(t *testing.T) {
	if response, err := http.Get("http://localhost:8080/v1/periodos/25"); err == nil {
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
