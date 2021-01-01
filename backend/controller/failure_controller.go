package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func getFailures(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	fromString := ctx.Query("from")
	toString := ctx.Query("to")

	from, err := time.Parse(time.RFC3339, fromString)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("date must be in format '%s'", time.RFC3339)))
		return
	}

	to, err := time.Parse(time.RFC3339, toString)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("date must be in format '%s'", time.RFC3339)))
		return
	}

	failures, err := service.GetFailures(ctx.Request.Context(), serviceId, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	failureVos := model.MapFailureEntitiesToVos(failures)
	ctx.JSON(http.StatusOK, model.FailureWrapperVo{Data: failureVos})
}

func getFailureCount(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	fromString := ctx.Query("from")
	toString := ctx.Query("to")

	from, err := time.Parse(time.RFC3339, fromString)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("date must be in format '%s'", time.RFC3339)))
		return
	}

	to, err := time.Parse(time.RFC3339, toString)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("date must be in format '%s'", time.RFC3339)))
		return
	}

	failureCount, err := service.GetFailuresCount(ctx.Request.Context(), serviceId, from, to)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	ctx.JSON(http.StatusOK, model.MapFailureCountEntityToVo(failureCount))
}
