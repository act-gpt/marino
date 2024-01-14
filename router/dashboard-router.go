package router

import (
	"github.com/act-gpt/marino/controller"
	"github.com/act-gpt/marino/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetDashboardRouter(router *gin.Engine) {
	router.GET("/link", controller.Redirect)
	router.GET("/doc", controller.GoDoc)

	openRouter := router.Group("/open")
	{
		chatRoute := openRouter.Group("/chat")
		{
			chatRoute.GET("/sign/:id", gzip.Gzip(gzip.DefaultCompression), controller.Sign)
			chatRoute.GET("/conversation/:id", gzip.Gzip(gzip.DefaultCompression), controller.ConversationHistories)
			chatRoute.Use(middleware.OpenAuth())
			{
				chatRoute.GET("/bot/:id", gzip.Gzip(gzip.DefaultCompression), controller.GetBot)
				chatRoute.POST("/query/:id", controller.Query)
				chatRoute.POST("/like/:id", controller.Like)
			}
			chatRoute.GET("/:id", gzip.Gzip(gzip.DefaultCompression), controller.SegmentDetail)
		}
	}

	apiRouter := router.Group("/dashboard")
	//apiRouter.Use(middleware.GlobalWebRateLimit(), gzip.Gzip(gzip.DefaultCompression))
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		apiRouter.GET("/status", middleware.CriticalRateLimit(), controller.GetStatus)
		apiRouter.POST("/verification", middleware.CriticalRateLimit(), controller.SendEmailVerification)
		apiRouter.POST("/reset_password", middleware.CriticalRateLimit(), controller.SendPasswordResetEmail)
		apiRouter.POST("/user/reset", middleware.CriticalRateLimit(), controller.ResetPassword)
		apiRouter.GET("/oauth/email/bind", middleware.CriticalRateLimit(), middleware.OwnerAuth(), controller.EmailBind)
		apiRouter.GET("/js/integrity", middleware.CriticalRateLimit(), controller.ShaIntegrity)
		apiRouter.GET("/models", middleware.CriticalRateLimit(), controller.Models)
		apiRouter.POST("/check", middleware.GlobalWebRateLimit(), controller.CheckEngine)

		meRoute := apiRouter.Group("/me")
		meRoute.Use(middleware.AdminAuth())
		{
			meRoute.GET("", controller.GetSelf)
			meRoute.PUT("", controller.UpdateSelf)
			meRoute.DELETE("", controller.DeleteSelf)
		}

		templateRoute := apiRouter.Group("/templates")
		{
			templateRoute.Use(middleware.AdminAuth())
			{
				templateRoute.GET("", controller.GetGetTemplate)
			}
		}

		userRoute := apiRouter.Group("/users")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)
			userRoute.GET("/logout", controller.Logout)

			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.AdminAuth())
			{
				selfRoute.GET("/find", controller.FindhUsers)
			}
			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.AdminAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/search", controller.SearchUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.PUT("/:id", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}

		botRoute := apiRouter.Group("/bots")
		{
			botRoute.Use(middleware.AdminAuth())
			{
				botRoute.GET("/", controller.GetBostByOwner)
				botRoute.GET("/admin", controller.GetBostByAdmin)
				botRoute.GET("/messages/:id", controller.GetMessagesByBot)
				botRoute.GET("/setting/:id", controller.GetSetting)
				botRoute.PUT("/setting/:id", controller.UpdateSetting)
				botRoute.GET("/:id", controller.GetBot)
				botRoute.POST("/", controller.CreateBot)
				botRoute.POST("/template", controller.CreateBot)
				botRoute.PUT("/:id", controller.UpdateBot)
			}
			botRoute.Use(middleware.OwnerAuth())
			{
				botRoute.DELETE("/:id", controller.DeleteBot)
			}
		}

		configRoute := apiRouter.Group("/config")
		{

			configRoute.POST("", controller.SetSystemConfig)
			configRoute.GET("", controller.GetSystemConfig)
			configRoute.Use(middleware.AdminAuth())
			{

			}
		}

		folderRoute := apiRouter.Group("/folders")
		{
			folderRoute.Use(middleware.AdminAuth())
			{
				folderRoute.GET("/bot/:id", controller.GetFolderByBot)
				folderRoute.GET("/:id", controller.GetFolder)
				folderRoute.POST("/", controller.CreateFolder)
				folderRoute.PUT("/:id", controller.UpdateFolder)
				folderRoute.DELETE("/:id", controller.DeleteFolder)
			}
		}

		knowledgeRoute := apiRouter.Group("/knowledges")
		{
			knowledgeRoute.Use(middleware.GlobalWebRateLimit(), middleware.AdminAuth())
			{
				knowledgeRoute.GET("/", controller.GetKnowledgesByPath)
				knowledgeRoute.GET("/search", controller.SearchKnowledges)
				knowledgeRoute.GET("/:id", controller.GetKnowledge)
				knowledgeRoute.POST("/", controller.CreateKnowledge)
				knowledgeRoute.POST("/upload", controller.UploadKnowledges)
				knowledgeRoute.PUT("/:id", controller.UpdateKnowledge)
				knowledgeRoute.DELETE("/", controller.BatchDeleteKnowledge)
				knowledgeRoute.DELETE("/:id", controller.DeleteKnowledge)
			}
		}

		orgRoute := apiRouter.Group("/orgs")
		{
			orgRoute.Use(middleware.GlobalWebRateLimit(), middleware.AdminAuth())
			{
				orgRoute.GET("/", controller.GetMyOrganization)
				orgRoute.POST("/", controller.CreateOrganization)
			}
			orgRoute.Use(middleware.GlobalWebRateLimit(), middleware.OwnerAuth())
			{
				orgRoute.PUT("/:id", controller.UpdateOrganization)
				orgRoute.DELETE("/:id", controller.DeleteOrganization)
			}
			orgRoute.Use(middleware.GlobalWebRateLimit(), middleware.RootAuth())
			{
				orgRoute.GET("/admin", controller.GetAllOrganization)
				orgRoute.GET("/:id", controller.GetOrganization)
			}

		}
	}
}
