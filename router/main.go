package router

import (
	"embed"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine, buildFS embed.FS) {
	SetDashboardRouter(router)
	SetApiRouter(router)
	setWebRouter(router, buildFS)
}
