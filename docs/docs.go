// Package docs provides OpenAPI documentation registration for the ShareToken blockchain.
package docs

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterOpenAPIService registers the OpenAPI service with the provided router.
// This is a placeholder implementation for the ShareToken blockchain.
// In a production environment, this would serve the generated OpenAPI specification.
func RegisterOpenAPIService(appName string, router *mux.Router) {
	// Serve a simple OpenAPI spec endpoint
	router.HandleFunc("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		// nolint:errcheck // OpenAPI spec is a constant, write will not fail
		w.Write([]byte(openAPISpec))
	}).Methods("GET")

	// Serve Swagger UI (optional, for development)
	router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		// nolint:errcheck // Swagger UI is a constant, write will not fail
		w.Write([]byte(swaggerUI))
	}).Methods("GET")
}

// OpenAPI specification placeholder
const openAPISpec = `
openapi: 3.0.0
info:
  title: ShareToken Blockchain API
  description: ShareToken blockchain HTTP API specification
  version: 1.0.0
  contact:
    name: ShareToken Team
servers:
  - url: http://localhost:1317
    description: Local development server
paths:
  /cosmos/bank/v1beta1/balances/{address}:
    get:
      summary: Get account balances
      tags:
        - Bank
      parameters:
        - name: address
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Successful response
`

// Swagger UI HTML placeholder
const swaggerUI = `<!DOCTYPE html>
<html>
<head>
    <title>ShareToken API - Swagger UI</title>
</head>
<body>
    <h1>ShareToken Blockchain API</h1>
    <p>View the OpenAPI specification at <a href="/openapi.yml">/openapi.yml</a></p>
</body>
</html>
`
