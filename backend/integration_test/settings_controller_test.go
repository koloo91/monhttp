package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)

func (suite *MonHttpTestSuite) TestIsSetupShouldBeFalse() {
	suite.resetSettings()

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/setup", nil)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), false, responseBody["isSetup"])
}

func (suite *MonHttpTestSuite) TestIsSetupShouldBeTrueAfterSetup() {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/setup", nil)

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), true, responseBody["isSetup"])
}

func (suite *MonHttpTestSuite) TestShouldSendSettingsIfNotSetup() {
	suite.resetSettings()

	requestBody, err := json.Marshal(map[string]interface{}{
		"databaseHost":     suite.databaseHost,
		"databasePort":     suite.databasePort,
		"databaseUser":     databaseUser,
		"databasePassword": databasePassword,
		"databaseName":     databaseName,
		"username":         user,
		"password":         password,
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/settings", bytes.NewBuffer(requestBody))

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), http.StatusNoContent, recorder.Code)
	assert.Equal(suite.T(), "", recorder.Body.String())

	configFile, err := os.Open("./config/config.env")
	assert.Nil(suite.T(), err)
	defer configFile.Close()

	contentBytes, err := ioutil.ReadAll(configFile)
	assert.Nil(suite.T(), err)

	content := string(contentBytes)

	assert.Contains(suite.T(), content, fmt.Sprintf("DATABASE_HOST=%s", suite.databaseHost))
	assert.Contains(suite.T(), content, fmt.Sprintf("DATABASE_PORT=%d", suite.databasePort))
	assert.Contains(suite.T(), content, fmt.Sprintf("DATABASE_USER=%s", databaseUser))
	assert.Contains(suite.T(), content, fmt.Sprintf("DATABASE_PASSWORD=%s", databasePassword))
	assert.Contains(suite.T(), content, fmt.Sprintf("DATABASE_NAME=%s", databaseName))
	assert.Contains(suite.T(), content, fmt.Sprintf("USERS=%s", userAndPassword))
}
