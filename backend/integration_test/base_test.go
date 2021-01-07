package integration_test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/monhttp/controller"
	"github.com/koloo91/monhttp/repository"
	"github.com/koloo91/monhttp/service"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
)

type MonHttpTestSuite struct {
	suite.Suite
	postgresContainer testcontainers.Container
	router            *gin.Engine
}

func (suite *MonHttpTestSuite) SetupSuite() {
	log.Println("Setup suite")

	container, host, port := startPostgresContainer(suite.T())
	suite.postgresContainer = container

	if err := service.LoadDatabase(host, port, "postgres", "postgres", "postgres"); err != nil {
		suite.Fail("Unable to load database: ", err)
	}

	repository.SetDatabase(service.GetDatabase())
	service.AddUser("admin", "admin")

	suite.router = controller.SetupRoutes()
}

func (suite *MonHttpTestSuite) SetupTest() {
	log.Println("Setup test")
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
	//defer postgresContainer.Terminate(ctx)
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
