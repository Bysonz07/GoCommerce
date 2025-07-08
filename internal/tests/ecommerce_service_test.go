package tests

import (
	"context"
	"testing"

	"ecommerce/api"
	"ecommerce/internal/service"
)

func TestAddToCart(t *testing.T) {
	svc := service.NewECommerceService(mockDB, mockRedis)

	req := &api.AddToCartRequest{
		UserId:    1,
		ProductId: 1,
		Quantity:  2,
	}

	resp, err := svc.AddToCart(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to add to cart: %v", err)
	}

	if resp.Message != "Product added to cart successfully" {
		t.Fatalf("Unexpected response: %v", resp.Message)
	}
}
