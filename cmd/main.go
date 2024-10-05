package main

import "eMobile/internal/app"

// @title           Music service
// @version         1.0
// @description     This is my server.

// @license.name  Apache helicopter
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app.Run()
}
