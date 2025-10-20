package v1_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func (ts *TestSuite) TestListPartSuccess() {
	partID1 := gofakeit.UUID()
	partID2 := gofakeit.UUID()

	ts.service.
		On("List", ts.ctx, mock.Anything).
		Return([]model.Part{{ID: partID1}, {ID: partID2}}, nil)

	resp, err := ts.api.ListParts(ts.ctx, &inventoryV1.ListPartsRequest{})

	ts.Require().NoError(err)
	ts.Require().Len(resp.Parts, 2)
	ts.Require().Equal(partID1, resp.Parts[0].Uuid)
	ts.Require().Equal(partID2, resp.Parts[1].Uuid)
}

func (ts *TestSuite) TestListPartGotError() {
	ts.service.
		On("List", ts.ctx, mock.Anything).
		Return(nil, gofakeit.Error())

	resp, err := ts.api.ListParts(ts.ctx, &inventoryV1.ListPartsRequest{})

	ts.Require().Error(err)
	ts.Require().Nil(resp)

	code, ok := status.FromError(err)

	ts.Require().True(ok)
	ts.Require().Equal(codes.Internal, code.Code())
}
