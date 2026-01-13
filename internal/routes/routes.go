package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/user/go-todo-api/docs"
	"github.com/user/go-todo-api/internal/handlers"
	"github.com/user/go-todo-api/internal/middleware"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Prometheus metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Custom logger middleware
	r.Use(middleware.LoggerMiddleware())

	// Rate limiting: 100 requests per minute per IP
	r.Use(middleware.RateLimitMiddleware(100, time.Minute))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	todoHandler := handlers.NewTodoHandler()

	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
		})

		// Swagger documentation
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("profile", authHandler.GetProfile)
			protected.POST("auth/logout", authHandler.Logout)

			// Todo routes
			todos := protected.Group("todos")
			{
				todos.GET("", todoHandler.GetAll)
				todos.GET("/:id", todoHandler.GetByID)
				todos.POST("", todoHandler.Create)
				todos.PUT("/:id", todoHandler.Update)
				todos.DELETE("/:id", todoHandler.Delete)
			}
		}
	}

	return r
}
