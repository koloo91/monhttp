package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type GetChecksQueryParameter struct {
	From     *time.Time `form:"from"`
	To       *time.Time `form:"to"`
	Interval *int       `form:"interval"`
}

func getChecks(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	var queryParameter GetChecksQueryParameter
	if err := ctx.ShouldBindQuery(&queryParameter); err != nil {
		log.Errorf("Unable to get query parameter: '%s'", err)
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	if queryParameter.From == nil {
		from := time.Now().Add(-24 * time.Hour)
		queryParameter.From = &from
	}

	if queryParameter.To == nil {
		to := time.Now()
		queryParameter.To = &to
	}

	if queryParameter.Interval == nil {
		interval := 300
		queryParameter.Interval = &interval
	}

	checks, err := service.GetChecks(ctx.Request.Context(), serviceId, *queryParameter.From, *queryParameter.To, *queryParameter.Interval)
	if err != nil {
		log.Errorf("Unable to get checks from database: '%s'", err)
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
		log.Errorf("Unable to get average values from database: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	averageVo := model.MapAverageEntityToVo(average)
	ctx.JSON(http.StatusOK, averageVo)
}

func getIsOnline(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	isOnline, err := service.GetIsOnline(ctx.Request.Context(), serviceId)
	if err != nil {
		log.Errorf("Unable to get is online value from database: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	ctx.JSON(http.StatusOK, model.IsOnlineVo{Online: isOnline})
}
