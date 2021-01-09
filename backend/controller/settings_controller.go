package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getSetup(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"isSetup": service.IsSetup()})
}

func postSettings(ctx *gin.Context) {
	if service.IsSetup() {
		log.Info("Is already setup, skipping")
		ctx.String(http.StatusNoContent, "")
		return
	}

	var settingsVo model.SettingsVo
	if err := ctx.ShouldBindJSON(&settingsVo); err != nil {
		log.Errorf("Unable to bind json body: '%s'", err)
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	if err := service.UpdateSettings(settingsVo); err != nil {
		log.Errorf("Unable to update settings: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	ctx.String(http.StatusNoContent, "")
}
