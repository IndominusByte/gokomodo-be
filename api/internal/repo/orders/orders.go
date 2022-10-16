package orders

import (
	"context"
	"fmt"

	ordersentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/orders"
	"github.com/creent-production/cdk-go/pagination"
	"github.com/jmoiron/sqlx"
)

type RepoOrders struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getOrderByDynamic": `SELECT transaction.orders.id, transaction.orders.buyer_id, buyer.email AS buyer_email, transaction.orders.seller_id, seller.email AS seller_email, transaction.orders.product_id, transaction.products.name AS product_name, transaction.orders.source_address, transaction.orders.destination_address, transaction.orders.qty, transaction.orders.price, transaction.orders.total_price, transaction.orders.status, transaction.orders.created_at, transaction.orders.updated_at FROM transaction.orders LEFT JOIN account.users AS buyer ON buyer.id = transaction.orders.buyer_id LEFT JOIN account.users AS seller ON seller.id = transaction.orders.seller_id LEFT JOIN transaction.products ON transaction.products.id = transaction.orders.product_id`,
}
var execs = map[string]string{
	"insertOrder": `INSERT INTO transaction.orders (buyer_id, seller_id, product_id, source_address, destination_address, qty, price, total_price) VALUES (:buyer_id, :seller_id, :product_id, :source_address, :destination_address, :qty, :price, :total_price) RETURNING id`,
}

func New(db *sqlx.DB) (*RepoOrders, error) {
	rp := &RepoOrders{
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
func (r *RepoOrders) Validate() error {
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

func (r *RepoOrders) GetProductById(ctx context.Context, id int) (*ordersentity.Product, error) {
	var t ordersentity.Product

	stmt, _ := r.db.PreparexContext(ctx, "SELECT transaction.products.id, transaction.products.price, transaction.products.seller_id, account.users.address AS seller_address FROM transaction.products LEFT JOIN account.users ON account.users.id = transaction.products.seller_id WHERE transaction.products.id = $1")
	return &t, stmt.GetContext(ctx, &t, id)
}

func (r *RepoOrders) GetOrderById(ctx context.Context, id, sellerId int) (*ordersentity.Order, error) {
	var t ordersentity.Order

	stmt, _ := r.db.PreparexContext(ctx, r.queries["getOrderByDynamic"]+` WHERE transaction.orders.id = $1 AND transaction.orders.seller_id = $2`)
	return &t, stmt.GetContext(ctx, &t, id, sellerId)
}

func (r *RepoOrders) Insert(ctx context.Context, payload *ordersentity.JsonCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertOrder"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}

func (r *RepoOrders) SetAcceptedOrder(ctx context.Context, id int) error {
	stmt, _ := r.db.PreparexContext(ctx, "UPDATE transaction.orders SET updated_at=CURRENT_TIMESTAMP, status='accepted' WHERE id = $1")
	_, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoOrders) GetOrderPaginate(ctx context.Context, payload *ordersentity.QueryOrderSchema) (*ordersentity.OrderPaginate, error) {
	var results ordersentity.OrderPaginate

	query := r.queries["getOrderByDynamic"] + ` WHERE 1=1`
	if len(payload.Q) > 0 {
		query += ` AND transaction.products.name LIKE '%'|| :q ||'%'`
	}
	switch payload.For {
	case "buyer":
		query += ` AND transaction.orders.buyer_id = :buyer_id`
	case "seller":
		query += ` AND transaction.orders.seller_id = :seller_id`
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
