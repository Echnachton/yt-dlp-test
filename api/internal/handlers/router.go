package handlers

import (
	"github.com/Echnachton/yt-dlp-test/internal/jobManager"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jbmgr *jobManager.JobManager) {
	apiV1Routes := router.Group("/api/v1")
	
	apiV1Routes.GET("/health", HealthHandler)
	apiV1Routes.GET("/health/:id", DynamicHealthFileHandler)

	apiV1Routes.POST("/file", PostFileHandler)
	
	apiV1Routes.POST("/download", PostDownloadHandler(jbmgr))
	apiV1Routes.GET("/download", GetDownloadsHandler)
	
	router.StaticFS("/web", gin.Dir("../web/static", true))
}
