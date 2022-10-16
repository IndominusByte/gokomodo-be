package auth

import (
	"context"

	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type RepoAuth struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getUserByDynamic": `SELECT id, email, password FROM account.users`,
}
var execs = map[string]string{}

func New(db *sqlx.DB) (*RepoAuth, error) {
	rp := &RepoAuth{
		db:      db,
		queries: queries,
		execs:   execs,
	}

	err := rp.Validate()
	if err != nil {
		return nil, err
	}

	return rp, nil
}

// Validate will validate sql query to db
func (r *RepoAuth) Validate() error {
	for _, q := range r.queries {
		_, err := r.db.PrepareNamedContext(context.Background(), q)
		if err != nil {
			return err
		}
	}

	for _, e := range r.execs {
		_, err := r.db.PrepareNamedContext(context.Background(), e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoAuth) IsPasswordSameAsHash(ctx context.Context, hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false
	}
	return true
}

func (r *RepoAuth) GetUserByEmail(ctx context.Context, email string) (*authentity.User, error) {
	var t authentity.User
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getUserByDynamic"]+" WHERE email = :email")

	return &t, stmt.GetContext(ctx, &t, authentity.User{Email: email})
}

func (r *RepoAuth) GetUserById(ctx context.Context, userId int) (*authentity.User, error) {
	var t authentity.User
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getUserByDynamic"]+" WHERE id = :id")

	return &t, stmt.GetContext(ctx, &t, authentity.User{Id: userId})
}
