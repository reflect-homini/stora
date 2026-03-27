package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/reflect-homini/stora/internal/adapters/http/handler"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/domain/appconstant"
)

func RegisterAPIRoutes(router *gin.Engine, handlers *handler.Handlers, authMiddleware gin.HandlerFunc) {
	apiRoutes := router.Group("/api")
	{
		v1 := apiRoutes.Group("/v1")
		{
			authRoutes := v1.Group("/auth")
			{
				authRoutes.POST("/register", handlers.Auth.HandleRegister())
				authRoutes.POST("/login", handlers.Auth.HandleInternalLogin())
				authRoutes.PUT("/refresh", handlers.Auth.HandleRefreshToken())
				authRoutes.GET(fmt.Sprintf("/:%s", string(appconstant.ContextProvider)), handlers.Auth.HandleOAuth2Login())
				authRoutes.GET(fmt.Sprintf("/:%s/callback", string(appconstant.ContextProvider)), handlers.Auth.HandleOAuth2Callback())
				authRoutes.GET("/verify-registration", handlers.Auth.HandleVerifyRegistration())
				authRoutes.POST("/password-reset", handlers.Auth.HandleSendPasswordReset())
				authRoutes.PATCH("/reset-password", handlers.Auth.HandleResetPassword())
			}

			protectedRoutes := v1.Group("/", authMiddleware)
			{
				protectedRoutes.DELETE("/auth/logout", handlers.Auth.HandleLogout())
				protectedRoutes.GET("/me", handlers.Auth.HandleMe())

				projectRoutes := protectedRoutes.Group("/projects")
				{
					projectRoutes.POST("", handlers.Project.HandleCreate())
					projectRoutes.GET("", handlers.Project.HandleGetAll())
					projectRoutes.GET("/:"+string(appconstant.ContextProjectID), handlers.Project.HandleGetByID())

					entryRoutes := projectRoutes.Group("/:" + string(appconstant.ContextProjectID) + "/entries")
					{
						entryRoutes.POST("", handlers.Project.HandleAddEntry())
						entryRoutes.GET("", handlers.Project.HandleGetEntriesAfter())
						entryRoutes.PUT("/:"+string(appconstant.ContextEntryID), handlers.Entry.HandleUpdateEntry())
						entryRoutes.DELETE("/:"+string(appconstant.ContextEntryID), handlers.Entry.HandleDeleteEntry())
					}
				}
			}

			if config.Global.App.Env == "debug" {
				v1.POST("/projects/:"+string(appconstant.ContextProjectID)+"/summaries", handlers.Project.HandleGenerateSummary())
				v1.POST("/projects/summaries", handlers.Project.HandleGenerateSummaries())
			}
		}
	}
}
