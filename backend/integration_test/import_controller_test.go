package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/koloo91/monhttp/service"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createMultipartFormBodyFromFile(fileName string, t *testing.T) (*bytes.Buffer, *multipart.Writer) {
	csvFile, err := os.Open(fileName)
	assert.Nil(t, err)

	fileContents, err := ioutil.ReadAll(csvFile)
	assert.Nil(t, err)

	csvFile.Close()

	requestBody := new(bytes.Buffer)
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", "file")
	assert.Nil(t, err)

	part.Write(fileContents)

	err = writer.Close()
	assert.Nil(t, err)

	return requestBody, writer
}

func (suite *MonHttpTestSuite) TestImportShouldReturnNotSetupIfNotSetup() {
	suite.resetSettings()

	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_one_ok.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	assert.Equal(suite.T(), "monhttp needs to be setup", responseBody["message"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnUnauthorizedWithoutCredentials() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_one_ok.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)
	assert.Equal(suite.T(), "invalid credentials", responseBody["message"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnOk() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_one_ok.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.SetBasicAuth(user, password)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 1, len(data))

	firsEntry := data[0].(map[string]interface{})
	assert.Equal(suite.T(), 1.0, firsEntry["rowNumber"])
	assert.Equal(suite.T(), "", firsEntry["error"])

	service := firsEntry["service"].(map[string]interface{})

	assert.NotEmpty(suite.T(), service["id"])
	assert.Equal(suite.T(), "Test Service", service["name"])
	assert.Equal(suite.T(), "HTTP", service["type"])
	assert.Equal(suite.T(), float64(30), service["intervalInSeconds"])
	assert.Equal(suite.T(), "http://google.com", service["endpoint"])
	assert.Equal(suite.T(), "GET", service["httpMethod"])
	assert.Equal(suite.T(), float64(60), service["requestTimeoutInSeconds"])
	assert.Equal(suite.T(), "", service["httpHeaders"])
	assert.Equal(suite.T(), "", service["httpBody"])
	assert.Equal(suite.T(), "", service["expectedHttpResponseBody"])
	assert.Equal(suite.T(), float64(200), service["expectedHttpStatusCode"])
	assert.Equal(suite.T(), true, service["followRedirects"])
	assert.Equal(suite.T(), true, service["verifySsl"])
	assert.Equal(suite.T(), true, service["enableNotifications"])
	assert.Equal(suite.T(), float64(2), service["notifyAfterNumberOfFailures"])
	assert.Equal(suite.T(), false, service["continuouslySendNotifications"])
	assert.Equal(suite.T(), []interface{}{"global", "telegram"}, service["notifiers"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnInvalidServiceType() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_invalid_service_type.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.SetBasicAuth(user, password)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 1, len(data))

	firsEntry := data[0].(map[string]interface{})
	assert.Equal(suite.T(), 1.0, firsEntry["rowNumber"])
	assert.Equal(suite.T(), service.ErrInvalidServiceType.Error(), firsEntry["error"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnInvalidHttpMethod() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_invalid_http_method.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.SetBasicAuth(user, password)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 1, len(data))

	firsEntry := data[0].(map[string]interface{})
	assert.Equal(suite.T(), 1.0, firsEntry["rowNumber"])
	assert.Equal(suite.T(), service.ErrInvalidHttpMethod.Error(), firsEntry["error"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnInvalidCheckIntervalInSeconds() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_invalid_check_interval_in_seconds.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.SetBasicAuth(user, password)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 1, len(data))

	firsEntry := data[0].(map[string]interface{})
	assert.Equal(suite.T(), 1.0, firsEntry["rowNumber"])
	assert.Equal(suite.T(), service.ErrInvalidIntervalInSeconds.Error(), firsEntry["error"])
}

func (suite *MonHttpTestSuite) TestImportShouldReturnInvalidRequestTimeoutInSeconds() {
	requestBody, multipartWriter := createMultipartFormBodyFromFile("files/csv/services_invalid_request_timeout_in_seconds.csv", suite.T())

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/import", requestBody)
	request.SetBasicAuth(user, password)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	data := responseBody["data"].([]interface{})
	assert.Equal(suite.T(), 1, len(data))

	firsEntry := data[0].(map[string]interface{})
	assert.Equal(suite.T(), 1.0, firsEntry["rowNumber"])
	assert.Equal(suite.T(), service.ErrInvalidRequestTimeoutInSeconds.Error(), firsEntry["error"])
}
