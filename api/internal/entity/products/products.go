package products

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type JsonCreateSchema struct {
	Name        string `json:"name" validate:"required,min=3,max=100" db:"name"`
	Description string `json:"description" validate:"required,min=5" db:"description"`
	Price       int    `json:"price" validate:"required,min=1" db:"price"`
	SellerId    int    `json:"-" db:"seller_id"`
}

type Product struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       int       `json:"price" db:"price"`
	SellerId    int       `json:"seller_id" db:"seller_id"`
	SellerName  string    `json:"seller_name" db:"seller_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductPaginate struct {
	Data      []Product  `json:"data"`
	Total     int        `json:"total"`
	NextNum   null.Int   `json:"next_num"`
	PrevNum   null.Int   `json:"prev_num"`
	Page      int        `json:"page"`
	IterPages []null.Int `json:"iter_pages"`
}

type QueryProductSchema struct {
	Page     int    `schema:"page" validate:"required,gte=1"`
	PerPage  int    `schema:"per_page" validate:"required,gte=1" db:"per_page"`
	Q        string `schema:"q" db:"q"`
	Offset   int    `schema:"-" db:"offset"`
	SellerId int    `schema:"-" db:"seller_id"`
}

type Exists struct {
	Exists bool `db:"exists"`
}
