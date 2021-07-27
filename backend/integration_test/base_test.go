package integration_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/controller"
	"github.com/koloo91/monhttp/repository"
	"github.com/koloo91/monhttp/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

const (
	databaseUser     = "postgres"
	databasePassword = "postgres"
	databaseName     = "postgres"

	databaseMigrationsDirectory = "../migrations"

	user     = "admin"
	password = "admin"

	userAndPassword = "admin:admin"
)

type MonHttpTestSuite struct {
	suite.Suite
	postgresContainer testcontainers.Container
	databaseHost      string
	databasePort      int
	router            *gin.Engine
}

func (suite *MonHttpTestSuite) SetupSuite() {
	log.Println("Setup suite")

	container, host, port := startPostgresContainer(suite.T())
	suite.postgresContainer = container
	suite.databaseHost = host
	suite.databasePort = port

	if err := service.LoadDatabase(host, port, databaseUser, databasePassword, databaseName, databaseMigrationsDirectory); err != nil {
		suite.Fail("Unable to load database: ", err)
	}

	repository.SetDatabase(service.GetDatabase())

	assert.Nil(suite.T(), service.AddUser("admin", "admin"))

	suite.router = controller.SetupRoutes()
}

func (suite *MonHttpTestSuite) SetupTest() {
	log.Println("Setup test")
	suite.setupSettings()
	suite.cleanupDatabase()
}

func (suite *MonHttpTestSuite) TearDownTest() {
	log.Println("Tear down test")

	os.Remove("./config/config.env")
	suite.resetSettings()
}

func (suite *MonHttpTestSuite) TearDownSuite() {
	log.Println(suite.postgresContainer.Terminate(context.Background()))
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(MonHttpTestSuite))
}

func startPostgresContainer(t *testing.T) (testcontainers.Container, string, int) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "postgres"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}

	ip, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Error(err)
	}

	return postgresContainer, ip, port.Int()
}

func (suite *MonHttpTestSuite) setupSettings() {
	assert.Nil(suite.T(), os.Setenv("DATABASE_HOST", suite.databaseHost))
	assert.Nil(suite.T(), os.Setenv("DATABASE_PORT", fmt.Sprintf("%d", suite.databasePort)))
	assert.Nil(suite.T(), os.Setenv("DATABASE_USER", databaseUser))
	assert.Nil(suite.T(), os.Setenv("DATABASE_PASSWORD", databasePassword))
	assert.Nil(suite.T(), os.Setenv("DATABASE_NAME", databaseName))

	assert.Nil(suite.T(), os.Setenv("USERS", userAndPassword))

	assert.Nil(suite.T(), os.Setenv("SCHEDULER_ENABLED", "false"))

	suite.loadTestConfig()
}

func (suite *MonHttpTestSuite) resetSettings() {
	assert.Nil(suite.T(), os.Setenv("DATABASE_HOST", ""))
	assert.Nil(suite.T(), os.Setenv("DATABASE_PORT", "0"))
	assert.Nil(suite.T(), os.Setenv("DATABASE_USER", ""))
	assert.Nil(suite.T(), os.Setenv("DATABASE_PASSWORD", ""))
	assert.Nil(suite.T(), os.Setenv("DATABASE_NAME", ""))

	assert.Nil(suite.T(), os.Setenv("USERS", ""))

	assert.Nil(suite.T(), os.Setenv("SCHEDULER_ENABLED", "false"))

	suite.loadTestConfig()
}

func (suite *MonHttpTestSuite) loadTestConfig() {
	assert.Nil(suite.T(), service.LoadConfig())
}

func (suite *MonHttpTestSuite) cleanupDatabase() {
	_, err := service.GetDatabase().Exec("DELETE FROM service;")
	if err != nil {
		suite.T().Fail()
	}
}
