package router

import (
	"embed"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/middleware"
	"github.com/act-gpt/marino/web"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
)

func setWebRouter(router *gin.Engine, buildFS embed.FS) {
	version := fmt.Sprintf("ACT GPT/%s (%s %s)", config.Version, runtime.GOARCH, runtime.GOOS)
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"version": version,
			"stats":   stats.Report(),
		})
	})
	index, _ := web.BuildFS.ReadFile("build/index.html")
	chat, _ := web.BuildFS.ReadFile("build/chat.html")
	router.Use(middleware.Cache())
	router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/locales") {
			method := middleware.CORS()
			method(c)
			return
		}
		c.Next()
	})

	router.GET("/js/embed/:id", func(c *gin.Context) {
		c.Status(200)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "application/javascript; charset=UTF-8")
		c.Writer.Write([]byte(common.EbbedFile()))
	})
	router.GET("/chat/:id", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", chat)
	})
	router.GET("/setup", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	router.GET("/", func(c *gin.Context) {
		if !system.Config.Initialled.Db {
			c.Redirect(http.StatusTemporaryRedirect, "/setup")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	router.Use(static.Serve("/", common.EmbedFolder(buildFS, "build")))
	router.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
}
