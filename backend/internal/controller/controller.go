package controller

import (
	"log/slog"
	"time"

	"github.com/christmas-fire/Bloomify/internal/metrics"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/christmas-fire/Bloomify/docs" // для генерации документации Swagger UI

	"github.com/gin-contrib/cors"
)

const defaultTimeout = 500 * time.Millisecond

type Handler struct {
	services  *service.Service
	validator *validator.Validate
	logger    *slog.Logger
	metrics   *metrics.Metrics
}

func NewHandler(services *service.Service, validator *validator.Validate, logger *slog.Logger, metrics *metrics.Metrics) *Handler {
	return &Handler{
		services:  services,
		validator: validator,
		logger:    logger,
		metrics:   metrics,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(h.LoggingMiddleware())
	router.Use(h.MetricsMiddleware())

	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.GET("/", h.getAllUsers)
				users.GET("/me", h.getMe)
				users.GET("/:id", h.getUserById)
				users.PATCH("/:id/username", h.updateUserUsername)
				users.PATCH("/:id/password", h.updateUserPassword)
				users.DELETE("/:id", h.deleteUser)
			}

			flowers := v1.Group("/flowers")
			{
				flowers.POST("/", h.createFlower)
				flowers.GET("/", h.getFlowers)
				flowers.GET("/:id", h.getFlowerById)
				flowers.PATCH("/:id", h.updateFlower)
				flowers.DELETE("/:id", h.deleteFlower)
			}
		}
	}

	return router
}
