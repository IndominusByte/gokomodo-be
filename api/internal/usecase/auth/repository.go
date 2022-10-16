package auth

import (
	"context"

	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
)

type authRepo interface {
	IsPasswordSameAsHash(ctx context.Context, hash, password []byte) bool
	GetUserByEmail(ctx context.Context, email string) (*authentity.User, error)
	GetUserById(ctx context.Context, userId int) (*authentity.User, error)
}
