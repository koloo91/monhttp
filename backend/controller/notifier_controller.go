package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getNotifiers(ctx *gin.Context) {
	notifiers := service.GetNotifiers()
	notifiersVo := model.MapNotifiersToVos(notifiers)
	ctx.JSON(http.StatusOK, model.NotifierWrapperVo{Data: notifiersVo})
}

func updateNotifier(ctx *gin.Context) {
	id := ctx.Param("id")

	var body map[string]interface{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	log.Println(id)
}
