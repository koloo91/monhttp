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

type GetFailuresQueryParameter struct {
	PageSize *int       `form:"pageSize" binding:"required"`
	Page     *int       `form:"page" binding:"required"`
	From     *time.Time `form:"from" binding:"required"`
	To       *time.Time `form:"to" binding:"required"`
}

func getFailures(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	var queryParameter GetFailuresQueryParameter
	if err := ctx.ShouldBindQuery(&queryParameter); err != nil {
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	failures, err := service.GetFailures(ctx.Request.Context(), serviceId, *queryParameter.From, *queryParameter.To, *queryParameter.PageSize, *queryParameter.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	failureCount, err := service.GetFailuresCount(ctx.Request.Context(), serviceId, *queryParameter.From, *queryParameter.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	failureVos := model.MapFailureEntitiesToVos(failures)
	ctx.JSON(http.StatusOK, model.FailureWrapperVo{
		Data:       failureVos,
		TotalCount: failureCount.Count,
		PageSize:   *queryParameter.PageSize,
		Page:       *queryParameter.Page,
	})
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
