package orders

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JsonCreateSchema struct {
	ProductId          int    `json:"product_id" validate:"required,min=1" db:"product_id"`
	Qty                int    `json:"qty" validate:"required,min=1" db:"qty"`
	BuyerId            int    `json:"-" db:"buyer_id"`
	SellerId           int    `json:"-" db:"seller_id"`
	SourceAddress      string `json:"-" db:"source_address"`
	DestinationAddress string `json:"-" db:"destination_address"`
	Price              int    `json:"-" db:"price"`
	TotalPrice         int    `json:"-" db:"total_price"`
}

type Product struct {
	Id            int    `db:"id"`
	Price         int    `db:"price"`
	SellerId      int    `db:"seller_id"`
	SellerAddress string `db:"seller_address"`
}

type QueryOrderSchema struct {
	Page     int    `schema:"page" validate:"required,gte=1"`
	PerPage  int    `schema:"per_page" validate:"required,gte=1" db:"per_page"`
	For      string `schema:"for" validate:"required,oneof=buyer seller"`
	Q        string `schema:"q" db:"q"`
	BuyerId  int    `schema:"-" db:"buyer_id"`
	SellerId int    `schema:"-" db:"seller_id"`
	Offset   int    `schema:"-" db:"offset"`
}

type Order struct {
	Id                 int       `json:"id" db:"id"`
	BuyerId            int       `json:"buyer_id" db:"buyer_id"`
	BuyerEmail         string    `json:"buyer_email" db:"buyer_email"`
	SellerId           int       `json:"seller_id" db:"seller_id"`
	SellerEmail        string    `json:"seller_email" db:"seller_email"`
	ProductId          int       `json:"product_id" db:"product_id"`
	ProductName        string    `json:"product_name" db:"product_name"`
	SourceAddress      string    `json:"source_address" db:"source_address"`
	DestinationAddress string    `json:"destination_address" db:"destination_address"`
	Qty                int       `json:"qty" db:"qty"`
	Price              int       `json:"price" db:"price"`
	TotalPrice         int       `json:"total_price" db:"total_price"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type OrderPaginate struct {
	Data      []Order    `json:"data"`
	Total     int        `json:"total"`
	NextNum   null.Int   `json:"next_num"`
	PrevNum   null.Int   `json:"prev_num"`
	Page      int        `json:"page"`
	IterPages []null.Int `json:"iter_pages"`
}
