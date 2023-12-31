package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"todo-app/pkg/service"
)

type Handler struct {
	log      *logrus.Logger
	services *service.Service
}

func NewHandler(log *logrus.Logger, services *service.Service) *Handler {
	return &Handler{
		log:      log,
		services: services,
	}
}

func (h *Handler) InitRouts() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
		auth.POST("refresh-tokens", h.refreshTokens)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}
	return router
}
