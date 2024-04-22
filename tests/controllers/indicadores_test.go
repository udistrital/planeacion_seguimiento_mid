package controllers

import (
	"net/http"
	"testing"
)

func TestConsultarIndicadores(t *testing.T) {
	if response, err := http.Get("http://localhost:8080/v1/indicadores/6201d43f25e40c205608b459"); err == nil {
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
