package config

import (
	"bytes"
	"cdn/common/constants"
	"cdn/common/utils"
	"fmt"
	"path/filepath"
	"text/template"
)

type OpenAPITemplate struct {
	Docs    []byte
	Redocly []byte
	Scalar  []byte
	Swagger []byte
}

type OpenAPIDocsData struct {
	BasePath string
}

var OpenAPITemplates = &OpenAPITemplate{}

func loadTemplate(templateFileName string, data any) (body []byte, err error) {
	tmpl, err := template.ParseFiles(
		filepath.Join(constants.ASSET_OPEN_API_PATH, templateFileName),
	)
	if err != nil {
		errMsg := "Error loading templates!"
		err = fmt.Errorf("%s: %s %w", errMsg, err.Error(), err)
		return
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		errMsg := "Error executing template!"
		err = fmt.Errorf("%s: %s %w", errMsg, err.Error(), err)
		return
	}

	return buf.Bytes(), nil
}

// Loads OpenAPI templates from a specified location resources.
func LoadOpenAPITemplates() (err error) {
	OpenAPITemplates.Docs, err = loadTemplate("/docs.html", &OpenAPIDocsData{
		BasePath: Env.ApiGroup,
	})
	if err != nil {
		return
	}

	// Redocly
	OpenAPITemplates.Redocly, err = utils.ReadFile(constants.ASSET_OPEN_API_PATH + "/redocly.html")
	if err != nil {
		return
	}

	// Scalar
	OpenAPITemplates.Scalar, err = utils.ReadFile(constants.ASSET_OPEN_API_PATH + "/scalar.html")
	if err != nil {
		return
	}

	// Swagger
	OpenAPITemplates.Swagger, err = utils.ReadFile(constants.ASSET_OPEN_API_PATH + "/swagger.html")
	if err != nil {
		return
	}

	return
}
