package v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/milkrage/microservices-course-homework/inventory/internal/api/inventory/v1"
	"github.com/milkrage/microservices-course-homework/inventory/internal/service/mocks"
)

type TestSuite struct {
	suite.Suite

	ctx     context.Context // nolint: containedctx
	service *mocks.PartService
	api     *v1.PartHandler
}

func (ts *TestSuite) SetupTest() {
	ts.ctx = context.Background()
	ts.service = mocks.NewPartService(ts.T())
	ts.api = v1.NewPartHandler(ts.service)
}

func (ts *TestSuite) TearDownTest() {
}

func TestInventoryAPI(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
