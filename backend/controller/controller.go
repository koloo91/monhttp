package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	"net/http"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")

	{
		apiGroup.POST("/settings", postSettings)
	}

	apiGroup.Use(isSetup())

	{
		apiGroup.POST("/services", postService)
		apiGroup.GET("/services", getServices)
		apiGroup.GET("/services/:id", getService)
		apiGroup.PUT("/services/:id", putService)
		apiGroup.DELETE("/services/:id", deleteService)
	}

	{
		apiGroup.GET("/services/:id/checks", getChecks)
		apiGroup.GET("/services/:id/average", getAverage)
		apiGroup.GET("/services/:id/online", getIsOnline)
	}

	{
		apiGroup.GET("/services/:id/failures", getFailures)
	}

	{
		apiGroup.GET("/notifiers", getNotifiers)
		apiGroup.PUT("/notifiers/:id", updateNotifier)
	}

	return router
}

func toApiError(err error) model.ApiErrorVo {
	return model.ApiErrorVo{
		Message: err.Error(),
	}
}

func isSetup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !service.IsSetup() {
			ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("monhttp needs to be setup")))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
