package service

import (
	"context"
	"fmt"
	"log"

	"ecommerce/api"
	"ecommerce/db/sqlc"
	"ecommerce/internal/email"

	"github.com/redis/go-redis/v8"
)

type ECommerceService struct {
	db    *sqlc.Queries
	redis *redis.Client
}

func NewECommerceService(db *sqlc.Queries, redis *redis.Client) *ECommerceService {
	return &ECommerceService{db: db, redis: redis}
}

func (s *ECommerceService) AddToCart(ctx context.Context, req *api.AddToCartRequest) (*api.AddToCartResponse, error) {
	err := s.db.AddToCart(ctx, sqlc.AddToCartParams{
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &api.AddToCartResponse{Message: "Product added to cart successfully"}, nil
}

func (s *ECommerceService) Checkout(ctx context.Context, req *api.CheckoutRequest) (*api.CheckoutResponse, error) {
	carts, err := s.db.GetCartByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, cart := range carts {
		total += float64(cart.Quantity) * cart.ProductPrice
	}

	orderID, err := s.db.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID: req.UserId,
		Total:  total,
		Status: "Pending",
	})
	if err != nil {
		return nil, err
	}

	// Push job to Redis worker
	go func() {
		err := s.redis.LPush(ctx, "order_jobs", orderID).Err()
		if err != nil {
			log.Printf("Failed to push job to Redis: %v", err)
		}
	}()

	return &api.CheckoutResponse{Message: "Order created successfully"}, nil
}

func (s *ECommerceService) MakePayment(ctx context.Context, req *api.MakePaymentRequest) (*api.MakePaymentResponse, error) {
	order, err := s.db.GetOrderById(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	if order.Status != "Pending" {
		return nil, fmt.Errorf("Order is not in a pending state")
	}

	// Simulate payment processing
	order.Status = "Paid"
	err = s.db.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     req.OrderId,
		Status: "Paid",
	})
	if err != nil {
		return nil, err
	}

	// Send email notification
	go email.SendOrderConfirmation(order.UserEmail, order.ID)

	return &api.MakePaymentResponse{Message: "Payment successful"}, nil
}
