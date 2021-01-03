package controller

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func postService(ctx *gin.Context) {
	var vo model.ServiceVo
	if err := ctx.ShouldBindJSON(&vo); err != nil {
		log.Errorf("Unable to bind json body: '%s'", err)
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	entity := model.MapServiceVoToEntity(vo)
	createdEntity, err := service.CreateService(ctx.Request.Context(), entity)
	if err != nil {
		log.Errorf("Unable to store service into database: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	createdVo := model.MapServiceEntityToVo(createdEntity)
	ctx.JSON(http.StatusCreated, createdVo)
}

func getServices(ctx *gin.Context) {
	services, err := service.GetServices(ctx.Request.Context())
	if err != nil {
		log.Errorf("Unable to get services from database: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	vos := model.MapServiceEntitiesToVos(services)
	ctx.JSON(http.StatusOK, model.ServiceWrapperVo{Data: vos})
}

func getService(ctx *gin.Context) {
	serviceId := ctx.Param("id")
	serviceEntity, err := service.GetServiceById(ctx.Request.Context(), serviceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Service with id '%s' not found", serviceId)
			ctx.JSON(http.StatusNotFound, toApiError(err))
			return
		}
		log.Errorf("Unable to get service from database: '%s'", err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	serviceVo := model.MapServiceEntityToVo(serviceEntity)
	ctx.JSON(http.StatusOK, serviceVo)
}

func putService(ctx *gin.Context) {
	serviceId := ctx.Param("id")

	var requestBody model.ServiceVo
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		log.Errorf("Unable to bind json body: '%s'", err)
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	serviceEntity := model.MapServiceVoToEntity(requestBody)
	serviceEntity, err := service.UpdateServiceById(ctx.Request.Context(), serviceId, serviceEntity)
	if err != nil {
		log.Errorf("Unable to update service in database with id '%s' - '%s'", serviceId, err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	serviceVo := model.MapServiceEntityToVo(serviceEntity)
	ctx.JSON(http.StatusOK, serviceVo)
}

func deleteService(ctx *gin.Context) {
	serviceId := ctx.Param("id")
	if err := service.DeleteServiceById(ctx.Request.Context(), serviceId); err != nil {
		log.Errorf("Unable to delete service from database with id '%s' - '%s'", serviceId, err)
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}
	ctx.JSON(http.StatusNoContent, "")
}
