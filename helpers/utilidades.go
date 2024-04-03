package helpers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func GuardarDocumento(documentos []interface{}) []interface{} {
	var respuestaDocumentos []interface{}

	for _, documento := range documentos {

		if documento.(map[string]interface{})["file"] != nil {
			documento := map[string]interface{}{
				"IdTipoDocumento": documento.(map[string]interface{})["IdTipoDocumento"],
				"nombre":          documento.(map[string]interface{})["nombre"],
				"metadatos":       documento.(map[string]interface{})["metadatos"],
				"descripcion":     documento.(map[string]interface{})["descripcion"],
				"file":            documento.(map[string]interface{})["file"],
			}

			var documentoAuxiliar []map[string]interface{}
			documentoAuxiliar = append(documentoAuxiliar, documento)
			documentoSubido, errorDocumento := RegistrarDocumento(documentoAuxiliar)

			if errorDocumento == nil {
				documentoTemporal := map[string]interface{}{
					"Nombre":        documentoSubido.(map[string]interface{})["Nombre"].(string),
					"Enlace":        documentoSubido.(map[string]interface{})["Enlace"],
					"Id":            documentoSubido.(map[string]interface{})["Id"],
					"TipoDocumento": documentoSubido.(map[string]interface{})["TipoDocumento"],
					"Activo":        documentoSubido.(map[string]interface{})["Activo"],
				}

				respuestaDocumentos = append(respuestaDocumentos, documentoTemporal)
			}
		}
	}
	return respuestaDocumentos
}

func RegistrarDocumento(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}
	errRegDoc := request.SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"/document/uploadAnyFormat", "POST", &resultadoRegistro, documento)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}

}
