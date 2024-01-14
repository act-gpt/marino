package router

import (
	"github.com/act-gpt/marino/controller"
	"github.com/act-gpt/marino/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {

	apiRouter := router.Group("/v1")
	{
		chatRoute := apiRouter.Group("/chat")
		{
			chatRoute.POST("/:id", middleware.TokenAuth(), controller.Query)
		}
	}
}
