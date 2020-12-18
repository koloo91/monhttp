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

func getChecks(ctx *gin.Context) {
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

	checks, err := service.GetChecks(ctx.Request.Context(), serviceId, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	checkVos := model.MapCheckEntitiesToVos(checks)
	ctx.JSON(http.StatusOK, model.CheckWrapperVo{Data: checkVos})
}

func getAverage(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	average, err := service.GetAverageValues(ctx.Request.Context(), serviceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	averageVo := model.MapAverageEntityToVo(average)
	ctx.JSON(http.StatusOK, averageVo)
}
