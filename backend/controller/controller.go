package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	"net/http"
	"strings"
)

var (
	users = make(map[string]string)
)

func SetupRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(static.Serve("/", static.LocalFile("./public", false)))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	})

	apiGroup := router.Group("/api")

	{
		apiGroup.GET("/alive", alive)
		apiGroup.GET("/setup", getSetup)
		apiGroup.POST("/settings", postSettings)
	}

	apiGroup.Use(isSetup())
	apiGroup.Use(basicAuth())

	service.SetOnAdminSetCallback(func(username, password string) {
		users[username] = password
	})

	users = service.LoadUsers()

	{
		apiGroup.GET("/login", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
	}

	{
		apiGroup.POST("/services", postService)
		apiGroup.GET("/services", getServices)
		apiGroup.GET("/services/:id", getService)
		apiGroup.PUT("/services/:id", putService)
		apiGroup.DELETE("/services/:id", deleteService)
	}

	{
		apiGroup.GET("/services/:id/checks", getChecks)
		apiGroup.GET("/services/:id/average", getAverage)
		apiGroup.GET("/services/:id/online", getIsOnline)
	}

	{
		apiGroup.GET("/services/:id/failures", getFailures)
		apiGroup.GET("/services/:id/failures/count", getFailureCount)
		apiGroup.GET("/services/:id/failures/countByDay", getFailuresGroupedByDay)
	}

	{
		apiGroup.GET("/notifiers", getNotifiers)
		apiGroup.PUT("/notifiers/:id", updateNotifier)
	}

	return router
}

func toApiError(err error) model.ApiErrorVo {
	return model.ApiErrorVo{
		Message: err.Error(),
	}
}

func alive(ctx *gin.Context) {
	ctx.String(http.StatusNoContent, "")
}

func isSetup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !service.IsSetup() {
			ctx.JSON(http.StatusBadRequest, toApiError(fmt.Errorf("monhttp needs to be setup")))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func basicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
			ctx.Abort()
			return
		}

		basicAuthValues := strings.Split(authorizationHeader, " ")
		if len(basicAuthValues) != 2 {
			ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
			ctx.Abort()
			return
		}

		base64Encoded := basicAuthValues[1]
		decodedBytes, err := base64.StdEncoding.DecodeString(base64Encoded)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
			ctx.Abort()
			return
		}

		decodedString := string(decodedBytes)
		usernameAndPassword := strings.Split(decodedString, ":")
		if len(usernameAndPassword) != 2 {
			ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
			ctx.Abort()
			return
		}

		if password, exists := users[usernameAndPassword[0]]; exists {
			if password != usernameAndPassword[1] {
				ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
				ctx.Abort()
				return
			}

			ctx.Next()
			return
		}

		ctx.JSON(http.StatusUnauthorized, model.ApiErrorVo{Message: "invalid credentials"})
		ctx.Abort()
		return
	}
}
