package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/service"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

func SetupRoutes() *gin.Engine {
	binding.Validator = new(defaultValidator)

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
		apiGroup.POST("/notifiers/:id/test/up", testNotifierUpTemplate)
		apiGroup.POST("/notifiers/:id/test/down", testNotifierDownTemplate)
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

		if password, exists := service.GetUsers()[usernameAndPassword[0]]; exists {
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

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			return name
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

type fieldError struct {
	err validator.FieldError
}

func (q fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + q.err.Field() + "'")
	sb.WriteString(", condition: " + q.err.ActualTag())

	// Print condition parameters, e.g. oneof=red blue -> { red blue }
	if q.err.Param() != "" {
		sb.WriteString(" { " + q.err.Param() + " }")
	}

	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}

	return sb.String()
}
