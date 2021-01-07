package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (suite *MonHttpTestSuite) TestExample() {
	suite.Equal(true, false)
}

func (suite *MonHttpTestSuite) TestCreateNewService() {
	requestBody, err := json.Marshal(map[string]interface{}{
		"name": "Foo",
	})
	assert.Nil(suite.T(), err)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/service", bytes.NewBuffer(requestBody))

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(200, recorder.Code)
	suite.Equal(suite.T(), "pong", recorder.Body.String())
}
