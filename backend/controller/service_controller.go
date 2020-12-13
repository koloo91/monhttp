package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	"net/http"
)

func postService(ctx *gin.Context) {
	var vo model.ServiceVo
	if err := ctx.ShouldBindJSON(&vo); err != nil {
		ctx.JSON(http.StatusBadRequest, toApiError(err))
		return
	}

	entity := model.MapServiceVoToEntity(vo)
	createdEntity, err := service.CreateService(ctx.Request.Context(), entity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	createdVo := model.MapServiceEntityToVo(createdEntity)
	ctx.JSON(http.StatusCreated, createdVo)
}

func getServices(ctx *gin.Context) {
	services, err := service.GetServices(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, toApiError(err))
		return
	}

	vos := model.MapServiceEntitiesToVos(services)
	ctx.JSON(http.StatusOK, model.ServiceWrapperVo{Data: vos})
}

func getService(ctx *gin.Context) {

}

func putService(ctx *gin.Context) {

}

func deleteService(ctx *gin.Context) {

}
