package inventory

import inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"

type Client struct {
	inventory inventoryV1.InventoryServiceClient
}

func NewInventoryClient(inventory inventoryV1.InventoryServiceClient) *Client {
	return &Client{
		inventory: inventory,
	}
}
