package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *MonHttpTestSuite) TestCreateNewServiceShouldReturnNotSetupIfNotSetup() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name": "Foo",
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(requestBody))

	suite.router.ServeHTTP(recorder, request)

	var responseBody map[string]interface{}
	assert.Nil(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &responseBody))

	assert.Equal(suite.T(), 400, recorder.Code)
	assert.Equal(suite.T(), "monhttp needs to be setup", responseBody["message"])
}
