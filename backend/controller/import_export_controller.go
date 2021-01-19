package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func importCsv(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Errorf("Unable to get uploaded file: '%s'", err)
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	fileHeader, err := file.Open()
	if err != nil {
		log.Errorf("Unable to open file: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	results, err := service.ImportCsvData(ctx.Request.Context(), fileHeader)
	if err != nil {
		log.Errorf("Unable to import data: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}
	resultVos := model.MapImportResultEntitiesToVos(results)

	ctx.JSON(http.StatusOK, gin.H{"data": resultVos})
}
