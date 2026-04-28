package api

const swaggerUIHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Agent Runtime API - Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html, body {
      margin: 0;
      min-height: 100%;
      background: #0b1118;
    }
    body {
      color: #e8f3ff;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    .swagger-ui {
      filter: invert(0.92) hue-rotate(180deg);
    }
    .swagger-ui .scheme-container,
    .swagger-ui .opblock,
    .swagger-ui .info,
    .swagger-ui .model-box,
    .swagger-ui table,
    .swagger-ui textarea,
    .swagger-ui input {
      box-shadow: none !important;
    }
    .swagger-ui .topbar {
      display: none;
    }
    #swagger-ui {
      max-width: 1440px;
      margin: 0 auto;
      padding: 12px 16px 28px;
    }
    .fallback {
      padding: 16px;
      color: #91a4ba;
    }
    .fallback a {
      color: #25d7cf;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <noscript><div class="fallback">Swagger UI requires JavaScript. Open <a href="/openapi.json">/openapi.json</a>.</div></noscript>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        displayOperationId: false,
        defaultModelsExpandDepth: 1,
        defaultModelExpandDepth: 2,
        docExpansion: 'list',
        persistAuthorization: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: 'StandaloneLayout'
      });
    };
  </script>
</body>
</html>`
