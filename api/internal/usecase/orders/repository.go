package orders

import (
	"context"

	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
	ordersentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/orders"
)

type ordersRepo interface {
	GetProductById(ctx context.Context, id int) (*ordersentity.Product, error)
	GetOrderById(ctx context.Context, id, sellerId int) (*ordersentity.Order, error)
	Insert(ctx context.Context, payload *ordersentity.JsonCreateSchema) int
	GetOrderPaginate(ctx context.Context, payload *ordersentity.QueryOrderSchema) (*ordersentity.OrderPaginate, error)
	SetAcceptedOrder(ctx context.Context, id int) error
}

type authRepo interface {
	GetUserById(ctx context.Context, userId int) (*authentity.User, error)
}
