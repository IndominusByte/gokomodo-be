package products

import (
	"context"

	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
	productsentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/products"
)

type productsRepo interface {
	IsDuplicate(ctx context.Context, name string, sellerId int) (*productsentity.Exists, error)
	Insert(ctx context.Context, payload *productsentity.JsonCreateSchema) int
	GetProductPaginate(ctx context.Context,
		payload *productsentity.QueryProductSchema, use string) (*productsentity.ProductPaginate, error)
}

type authRepo interface {
	GetUserById(ctx context.Context, userId int) (*authentity.User, error)
}
