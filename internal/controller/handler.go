package controller

import (
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/christmas-fire/Bloomify/docs"

	"github.com/gin-contrib/cors"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

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

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
		auth.POST("refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.GET("/", h.getAllUsers)
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
