package v1_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/milkrage/microservices-course-homework/inventory/internal/api/inventory/v1"
	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	"github.com/milkrage/microservices-course-homework/inventory/internal/service/mocks"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func (ts *TestSuite) TestGetPartSuccess() {
	partID := gofakeit.UUID()

	ts.service.
		On("Get", ts.ctx, partID).
		Return(model.Part{ID: partID}, nil)

	part, err := ts.api.GetPart(ts.ctx, &inventoryV1.GetPartRequest{Uuid: partID})

	ts.Require().NoError(err)
	ts.Require().Equal(partID, part.Part.Uuid)
}

func (ts *TestSuite) TestGetPartGotError() {
	partID := gofakeit.UUID()

	testCases := []struct {
		name         string
		serviceError error
		apiRespCode  codes.Code
	}{
		{name: "not found err", serviceError: model.ErrPartNotFound, apiRespCode: codes.NotFound},
		{name: "internal err", serviceError: gofakeit.Error(), apiRespCode: codes.Internal},
	}

	for _, test := range testCases {
		ts.Run(test.name, func() {
			ts.service = mocks.NewPartService(ts.T())
			ts.api = v1.NewPartHandler(ts.service)

			ts.service.
				On("Get", ts.ctx, partID).
				Return(model.Part{}, test.serviceError)

			part, err := ts.api.GetPart(ts.ctx, &inventoryV1.GetPartRequest{Uuid: partID})

			ts.Require().Nil(part)
			ts.Require().Error(err)

			code, ok := status.FromError(err)
			ts.Require().Equal(true, ok)
			ts.Require().Equal(test.apiRespCode, code.Code())
		})
	}
}
