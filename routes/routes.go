package routes

import (
	"todos/handlers"
	"todos/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// Swagger docs route with dark theme and custom UI
	r.GET("/swagger/*any", func(c *gin.Context) {
		if c.Param("any") == "/index.html" {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html { box-sizing: border-box; overflow: hidden; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background: #0d1117; display: flex; height: 100vh; overflow: hidden; }
    #sidebar { width: 280px; background: #161b22; border-right: 1px solid #30363d; overflow-y: auto; padding: 20px 0; flex-shrink: 0; }
    #sidebar .sidebar-title { color: #58a6ff; font-size: 1.4rem; font-weight: bold; padding: 0 20px 20px; border-bottom: 1px solid #30363d; margin-bottom: 20px; }
    #sidebar .nav-item { padding: 10px 20px; color: #8b949e; cursor: pointer; transition: all 0.2s; font-size: 0.95rem; }
    #sidebar .nav-item:hover { background: #21262d; color: #58a6ff; }
    #sidebar .nav-item.active { background: #1f242c; color: #58a6ff; border-left: 4px solid #58a6ff; }
    #content { flex-grow: 1; overflow-y: auto; padding: 0 40px; }
    .swagger-ui { background-color: #0d1117; color: #c9d1d9; }
    .swagger-ui .topbar { display: none; }
    .swagger-ui .info { background-color: #0d1117; color: #c9d1d9; padding: 20px 0; border-bottom: 1px solid #30363d; margin-bottom: 20px; }
    .swagger-ui .info .title { color: #58a6ff; }
    .swagger-ui .scheme-container { background-color: #0d1117; color: #c9d1d9; border-top: 1px solid #30363d; box-shadow: none; padding: 10px 0; margin-bottom: 20px; }
    .swagger-ui select { background-color: #21262d; color: #c9d1d9; border: 1px solid #30363d; }
    .swagger-ui .opblock { border-radius: 8px; border: 1px solid #30363d; background: #161b22; margin-bottom: 15px; }
    .swagger-ui .opblock .opblock-summary { border-bottom: 1px solid #30363d; }
    .swagger-ui .opblock-tag { font-size: 1.2rem; border-bottom: 1px solid #30363d; margin-bottom: 10px; color: #eee; padding: 10px 0; }
    .swagger-ui .opblock .opblock-summary-method { border-radius: 6px; text-shadow: none; }
    .swagger-ui .opblock-get { border-color: #2ea043; background: rgba(46, 160, 67, 0.1); }
    .swagger-ui .opblock-get .opblock-summary-method { background: #2ea043; }
    .swagger-ui .opblock-post { border-color: #238636; background: rgba(35, 134, 54, 0.1); }
    .swagger-ui .opblock-post .opblock-summary-method { background: #238636; }
    .swagger-ui .opblock-put { border-color: #d29922; background: rgba(210, 153, 34, 0.1); }
    .swagger-ui .opblock-put .opblock-summary-method { background: #d29922; }
    .swagger-ui .opblock-delete { border-color: #f85149; background: rgba(248, 81, 73, 0.1); }
    .swagger-ui .opblock-delete .opblock-summary-method { background: #f85149; }
    .swagger-ui .btn.authorize { color: #2ea043; border-color: #2ea043; background-color: transparent; }
    .swagger-ui .btn.authorize svg { fill: #2ea043; }
    .swagger-ui .btn { color: #c9d1d9; border-color: #30363d; background: #21262d; }
    .swagger-ui section.models { border: 1px solid #30363d; border-radius: 8px; margin-top: 20px; }
    .swagger-ui section.models h4 { color: #8b949e; }
    .swagger-ui .model-box { background: #161b22; }
    .swagger-ui .parameter__name, .swagger-ui .parameter__type { color: #8b949e; }
    .swagger-ui .opblock-description-wrapper p, .swagger-ui .opblock-external-docs-wrapper p, .swagger-ui .opblock-title_normal p { color: #c9d1d9; }
    .swagger-ui table thead tr td, .swagger-ui table thead tr th { color: #8b949e; border-bottom: 1px solid #30363d; }
    .swagger-ui .response-col_status, .swagger-ui .response-col_description { color: #c9d1d9; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; }
  </style>
</head>
<body>
  <div id="sidebar">
    <div class="sidebar-title">Todos API</div>
    <div id="sidebar-nav"></div>
  </div>
  <div id="content">
    <div id="swagger-ui"></div>
  </div>
  <script src="./swagger-ui-bundle.js"> </script>
  <script src="./swagger-ui-standalone-preset.js"> </script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        persistAuthorization: true,
        docExpansion: "list",
        defaultModelsExpandDepth: -1,
        onComplete: function() {
          // Populate sidebar after Swagger UI loads
          setTimeout(() => {
            const sidebarNav = document.getElementById('sidebar-nav');
            const tags = document.querySelectorAll('.opblock-tag');
            tags.forEach((tag, index) => {
              const tagName = tag.innerText.split('\n')[0].trim();
              const navItem = document.createElement('div');
              navItem.className = 'nav-item';
              navItem.innerText = tagName;
              navItem.onclick = () => {
                tag.scrollIntoView({ behavior: 'smooth' });
                document.querySelectorAll('.nav-item').forEach(item => item.classList.remove('active'));
                navItem.classList.add('active');
              };
              sidebarNav.appendChild(navItem);
            });
          }, 1000);
        }
      })
      window.ui = ui
    }
  </script>
</body>
</html>
			`)
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			func(c *ginSwagger.Config) {
				c.InstanceName = "swagger"
				c.URL = "/swagger/doc.json"
			},
		)(c)
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Todos API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/profile", middleware.AuthRequired(), handlers.GetProfile)
		}

		// Todo routes (protected)
		todos := api.Group("/todos")
		todos.Use(middleware.AuthRequired())
		{
			todos.GET("", handlers.GetTodos)
			todos.GET("/:id", handlers.GetTodo)
			todos.POST("", handlers.CreateTodo)
			todos.PUT("/:id", handlers.UpdateTodo)
			todos.DELETE("/:id", handlers.DeleteTodo)
		}
	}
}
