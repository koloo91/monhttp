package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createMultipartFormBody(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
	csvFile, err := os.Open("files/services.csv")
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

	requestBody, multipartWriter := createMultipartFormBody(suite.T())

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
	requestBody, multipartWriter := createMultipartFormBody(suite.T())

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
	requestBody, multipartWriter := createMultipartFormBody(suite.T())

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
