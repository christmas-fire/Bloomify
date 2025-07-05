package controller

import (
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/metrics"
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/christmas-fire/Bloomify/docs" // для генерации онлайн-документации Swagger UI

	"github.com/gin-contrib/cors"
)

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

	// Глобальный обработчик OPTIONS для всех путей
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
				flowers.GET("/", h.getAllFlowers)
				flowers.GET("/:id", h.getFlowerById)
				flowers.GET("/name", h.getFlowersByName)
				flowers.GET("/description", h.getFlowersByDescription)
				flowers.GET("/price", h.getFlowersByPrice)
				flowers.GET("/stock", h.getFlowersByStock)
				flowers.PATCH("/:id/name", h.updateFlowerName)
				flowers.PATCH("/:id/description", h.updateFlowerDescription)
				flowers.PATCH("/:id/price", h.updateFlowerPrice)
				flowers.PATCH("/:id/stock", h.updateFlowerStock)
				flowers.DELETE("/:id", h.deleteFlower)
			}

			orders := v1.Group("/orders")
			{
				orders.POST("/", h.createOrder)
				orders.GET("/", h.getAllOrders)
				orders.GET("/:id", h.getOrderById)
				orders.GET("/user_id", h.getOrdersByUserId)
				orders.PUT("/:id", h.updateOrder)
				orders.PATCH("/:id/flower_id", h.updateOrderFlowerId)
				orders.PATCH("/:id/quantity", h.updateOrderQuantity)
				orders.DELETE("/:id", h.deleteOrder)
				orders.DELETE("/flower/:flower_id/", h.removeFlowerFromOrder)
				orders.PATCH("/flower/:flower_id/increment/", h.incrementFlowerQuantity)
				orders.PATCH("/flower/:flower_id/decrement/", h.decrementFlowerQuantity)
				orders.DELETE("/active", h.deleteActiveOrder)
			}

			order_flowers := v1.Group("/order_flowers")
			{
				order_flowers.GET("/", h.getAllOrderFlowers)
				order_flowers.GET("/:id", h.getOrderFlowersByOrderId)
			}
		}
	}

	return router
}
