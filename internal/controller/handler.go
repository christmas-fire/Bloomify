package controller

import (
	"github.com/christmas-fire/Bloomify/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/christmas-fire/Bloomify/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
				users.GET("/:id", h.getUserById)
				users.PATCH("/:id/change-username", h.updateUserUsername)
				users.PATCH("/:id/change-password", h.updateUserPassword)
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
				flowers.PATCH("/:id/change-name", h.updateFlowerName)
				flowers.PATCH("/:id/change-description", h.updateFlowerDescription)
				flowers.PATCH("/:id/change-price", h.updateFlowerPrice)
				flowers.PATCH("/:id/change-stock", h.updateFlowerStock)
				flowers.DELETE("/:id", h.deleteFlower)

			}

			orders := v1.Group("/orders")
			{
				orders.POST("/", h.createOrder)
				orders.GET("/", h.getAllOrders)
				orders.GET("/:id", h.getOrderById)
				orders.GET("/user_id", h.getOrdersByUserId)
			}
		}
	}

	return router
}
