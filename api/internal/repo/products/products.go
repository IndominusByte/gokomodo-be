package products

import (
	"context"
	"fmt"

	productsentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/products"
	"github.com/creent-production/cdk-go/pagination"
	"github.com/jmoiron/sqlx"
)

type RepoProducts struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getProductByDynamic": `SELECT transaction.products.id, transaction.products.name, transaction.products.description, transaction.products.price, 
transaction.products.seller_id, account.users.name AS seller_name, transaction.products.created_at, transaction.products.updated_at
FROM transaction.products LEFT JOIN account.users ON account.users.id = transaction.products.seller_id`,
}
var execs = map[string]string{
	"insertProduct": `INSERT INTO transaction.products (name, description, price, seller_id) VALUES (:name, :description, :price, :seller_id) RETURNING id`,
}

func New(db *sqlx.DB) (*RepoProducts, error) {
	rp := &RepoProducts{
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
func (r *RepoProducts) Validate() error {
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

func (r *RepoProducts) IsDuplicate(ctx context.Context, name string, sellerId int) (*productsentity.Exists, error) {
	var t productsentity.Exists

	stmt, _ := r.db.PreparexContext(ctx, "SELECT EXISTS (SELECT * FROM transaction.products WHERE name = $1 AND seller_id = $2)")
	return &t, stmt.GetContext(ctx, &t, name, sellerId)
}

func (r *RepoProducts) Insert(ctx context.Context, payload *productsentity.JsonCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertProduct"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}

func (r *RepoProducts) GetProductPaginate(ctx context.Context,
	payload *productsentity.QueryProductSchema, use string) (*productsentity.ProductPaginate, error) {

	var results productsentity.ProductPaginate

	query := r.queries["getProductByDynamic"] + ` WHERE 1=1`
	if len(payload.Q) > 0 {
		query += ` AND transaction.products.name LIKE '%'|| :q ||'%'`
	}
	if use == "seller" {
		query += ` AND transaction.products.seller_id = :seller_id`
	}

	// pagination
	var count struct{ Total int }
	stmt_count, _ := r.db.PrepareNamedContext(ctx, fmt.Sprintf("SELECT count(*) AS total FROM (%s) AS anon_1", query))
	err := stmt_count.GetContext(ctx, &count, payload)
	if err != nil {
		return &results, err
	}
	payload.Offset = (payload.Page - 1) * payload.PerPage

	// results
	query += ` LIMIT :per_page OFFSET :offset`
	stmt, _ := r.db.PrepareNamedContext(ctx, query)
	err = stmt.SelectContext(ctx, &results.Data, payload)
	if err != nil {
		return &results, err
	}

	paginate := pagination.Paginate{Page: payload.Page, PerPage: payload.PerPage, Total: count.Total}
	results.Total = paginate.Total
	results.NextNum = paginate.NextNum()
	results.PrevNum = paginate.PrevNum()
	results.Page = paginate.Page
	results.IterPages = paginate.IterPages()

	return &results, nil
}
