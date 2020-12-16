package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")

	{
		apiGroup.POST("/services", postService)
		apiGroup.GET("/services", getServices)
		apiGroup.GET("/services/:id", getService)
		apiGroup.PUT("/services/:id", putService)
		apiGroup.DELETE("/services/:id", deleteService)
	}

	{
		apiGroup.GET("/services/:id/checks", getChecks)
	}

	return router
}

func toApiError(err error) model.ApiErrorVo {
	return model.ApiErrorVo{
		Message: err.Error(),
	}
}
