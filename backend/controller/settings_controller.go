package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	"net/http"
)

func getSetup(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"isSetup": service.IsSetup()})
}

func postSettings(ctx *gin.Context) {
	if service.IsSetup() {
		ctx.JSON(http.StatusOK, "")
		return
	}

	var settingsVo model.SettingsVo
	if err := ctx.ShouldBindJSON(&settingsVo); err != nil {
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	if err := service.UpdateSettings(settingsVo); err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	ctx.JSON(http.StatusOK, "")
}
