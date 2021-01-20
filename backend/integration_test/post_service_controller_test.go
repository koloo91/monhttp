package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnNotSetupIfNotSetup() {
	suite.resetSettings()

	requestBody, err := json.Marshal(map[string]interface{}{
		"name": "Foo",
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), "monhttp needs to be setup", responseBody["message"])
}

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnUnauthorizedWithoutCredentials() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name": "Foo",
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)
	assert.Equal(suite.T(), "invalid credentials", responseBody["message"])
}

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnCreated() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":                          "MyService",
		"type":                          "HTTP",
		"intervalInSeconds":             30,
		"endpoint":                      "http://localhost",
		"httpMethod":                    "GET",
		"requestTimeoutInSeconds":       60,
		"httpHeaders":                   "httpHeaders",
		"httpBody":                      "httpBody",
		"expectedHttpResponseBody":      "expectedHttpResponseBody",
		"expectedHttpStatusCode":        200,
		"followRedirects":               true,
		"verifySsl":                     true,
		"enableNotifications":           true,
		"notifyAfterNumberOfFailures":   2,
		"continuouslySendNotifications": true,
		"notifiers":                     []string{"global"},
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))
	request.SetBasicAuth(user, password)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)

	assert.NotEmpty(suite.T(), responseBody["id"])
	assert.Equal(suite.T(), "MyService", responseBody["name"])
	assert.Equal(suite.T(), "HTTP", responseBody["type"])
	assert.Equal(suite.T(), float64(30), responseBody["intervalInSeconds"])
	assert.Equal(suite.T(), "http://localhost", responseBody["endpoint"])
	assert.Equal(suite.T(), "GET", responseBody["httpMethod"])
	assert.Equal(suite.T(), float64(60), responseBody["requestTimeoutInSeconds"])
	assert.Equal(suite.T(), "httpHeaders", responseBody["httpHeaders"])
	assert.Equal(suite.T(), "httpBody", responseBody["httpBody"])
	assert.Equal(suite.T(), "expectedHttpResponseBody", responseBody["expectedHttpResponseBody"])
	assert.Equal(suite.T(), float64(200), responseBody["expectedHttpStatusCode"])
	assert.Equal(suite.T(), true, responseBody["followRedirects"])
	assert.Equal(suite.T(), true, responseBody["verifySsl"])
	assert.Equal(suite.T(), true, responseBody["enableNotifications"])
	assert.Equal(suite.T(), float64(2), responseBody["notifyAfterNumberOfFailures"])
	assert.Equal(suite.T(), true, responseBody["continuouslySendNotifications"])
	assert.Equal(suite.T(), []interface{}{"global"}, responseBody["notifiers"])
}

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnErrorIfWrongType() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":                          "MyService",
		"type":                          "HTTPs",
		"intervalInSeconds":             30,
		"endpoint":                      "http://localhost",
		"httpMethod":                    "GET",
		"requestTimeoutInSeconds":       60,
		"httpHeaders":                   "httpHeaders",
		"httpBody":                      "httpBody",
		"expectedHttpResponseBody":      "expectedHttpResponseBody",
		"expectedHttpStatusCode":        200,
		"followRedirects":               true,
		"verifySsl":                     true,
		"enableNotifications":           true,
		"notifyAfterNumberOfFailures":   2,
		"continuouslySendNotifications": true,
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))
	request.SetBasicAuth(user, password)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), "Key: 'ServiceVo.Type' Error:Field validation for 'Type' failed on the 'oneof' tag", responseBody["message"])
}

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnErrorIfInvalidIntervalInSeconds() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":                          "MyService",
		"type":                          "HTTP",
		"intervalInSeconds":             29,
		"endpoint":                      "http://localhost",
		"httpMethod":                    "GET",
		"requestTimeoutInSeconds":       60,
		"httpHeaders":                   "httpHeaders",
		"httpBody":                      "httpBody",
		"expectedHttpResponseBody":      "expectedHttpResponseBody",
		"expectedHttpStatusCode":        200,
		"followRedirects":               true,
		"verifySsl":                     true,
		"enableNotifications":           true,
		"notifyAfterNumberOfFailures":   2,
		"continuouslySendNotifications": true,
		"notifiers":                     []string{},
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))
	request.SetBasicAuth(user, password)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), "Key: 'ServiceVo.IntervalInSeconds' Error:Field validation for 'IntervalInSeconds' failed on the 'min' tag", responseBody["message"])
}
