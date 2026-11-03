package docs

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed openapi.yaml
var OpenAPISpec []byte

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Go Template API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/docs/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`

func RegisterRoutes(ec *echo.Echo) {
	ec.GET("/docs/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/yaml", OpenAPISpec)
	})
	ec.GET("/swagger", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerUIHTML)
	})
}

