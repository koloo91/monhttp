package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/koloo91/monhttp/repository"
	"github.com/koloo91/monhttp/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *MonHttpTestSuite) createNewService() (string, error) {
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
	if err != nil {
		return "", err
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/services", bytes.NewBuffer(requestBody))
	request.SetBasicAuth(user, password)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &responseBody); err != nil {
		return "", nil
	}

	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)

	return responseBody["id"].(string), nil
}

func (suite *MonHttpTestSuite) TestServiceCheckShouldBeFailureIfHostNotReachable() {
	serviceId, err := suite.createNewService()
	assert.Nil(suite.T(), err)

	job, err := repository.GetJobByServiceId(context.Background(), serviceId)
	assert.Nil(suite.T(), err)

	service.ProcessService(0, job.Id)

	// get checks
	checksRecorder := httptest.NewRecorder()
	checksRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/services/%s/checks", serviceId), nil)
	checksRequest.SetBasicAuth(user, password)

	suite.router.ServeHTTP(checksRecorder, checksRequest)

	var checksResponseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(checksRecorder.Body.Bytes(), &checksResponseBody))

	assert.Equal(suite.T(), http.StatusOK, checksRecorder.Code)

	checksData := checksResponseBody["data"].([]interface{})

	firstChecksEntry := checksData[0].(map[string]interface{})
	assert.NotEmpty(suite.T(), firstChecksEntry["id"])
	assert.Equal(suite.T(), serviceId, firstChecksEntry["serviceId"])
	assert.Equal(suite.T(), 0.0, firstChecksEntry["latencyInMs"])
	assert.Equal(suite.T(), true, firstChecksEntry["isFailure"])
	assert.NotEmpty(suite.T(), firstChecksEntry["createdAt"])

	// get failures
	failuresRecorder := httptest.NewRecorder()
	failuresRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/services/%s/failures", serviceId), nil)
	failuresRequest.SetBasicAuth(user, password)

	suite.router.ServeHTTP(failuresRecorder, failuresRequest)

	var failuresResponseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(failuresRecorder.Body.Bytes(), &failuresResponseBody))

	assert.Equal(suite.T(), http.StatusOK, failuresRecorder.Code)
	// TODO: assert
}
