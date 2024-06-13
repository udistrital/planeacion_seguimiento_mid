package services

import (
	"encoding/json"
	"fmt"
)

func SolicitudReformulacion(requestBody []byte) (interface{}, error) {
	var body map[string]interface{}

	json.Unmarshal(requestBody, &body)
	fmt.Println(requestBody)

	return requestBody, nil
}
