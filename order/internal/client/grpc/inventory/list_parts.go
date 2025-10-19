package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/milkrage/microservices-course-homework/order/internal/client/converter"
	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (c *Client) ListParts(ctx context.Context, filter model.PartFilter) (model.ListParts, error) {
	req := converter.PartFilterToListPartsRequest(filter)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, err := c.inventory.ListParts(ctxWithTimeout, req)
	if err != nil {
		return model.ListParts{}, fmt.Errorf("%w: %w", model.ErrInternal, err)
	}

	listParts := converter.ListPartsResponseToModel(resp)

	return listParts, nil
}
