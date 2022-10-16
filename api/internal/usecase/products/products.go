package products

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	productsentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/products"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
	"github.com/go-chi/jwtauth"
)

type ProductsUsecase struct {
	productsRepo productsRepo
	authRepo     authRepo
}

func NewProductsUsecase(productRepo productsRepo, authRepo authRepo) *ProductsUsecase {
	return &ProductsUsecase{
		productsRepo: productRepo,
		authRepo:     authRepo,
	}
}

func (uc *ProductsUsecase) Create(ctx context.Context, rw http.ResponseWriter, payload *productsentity.JsonCreateSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	_, claims, _ := jwtauth.FromContext(ctx)
	sub, _ := strconv.Atoi(claims["sub"].(string))

	user, err := uc.authRepo.GetUserById(ctx, sub)
	if err != nil {
		response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
			constant.Header: constant.UserNotFound,
		})
		return
	}

	// check duplicate product
	if d, _ := uc.productsRepo.IsDuplicate(ctx, payload.Name, user.Id); d.Exists {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			constant.App: fmt.Sprintf(constant.AlreadyTaken, "product"),
		})
		return
	}

	// insert into db
	payload.SellerId = user.Id
	uc.productsRepo.Insert(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		constant.App: "Successfully add a new product.",
	})
}

func (uc *ProductsUsecase) MyProduct(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryProductSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	_, claims, _ := jwtauth.FromContext(ctx)
	sub, _ := strconv.Atoi(claims["sub"].(string))

	user, err := uc.authRepo.GetUserById(ctx, sub)
	if err != nil {
		response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
			constant.Header: constant.UserNotFound,
		})
		return
	}

	payload.SellerId = user.Id
	results, _ := uc.productsRepo.GetProductPaginate(ctx, payload, "seller")

	response.WriteJSONResponse(rw, 200, results, nil)
}

func (uc *ProductsUsecase) GetAll(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryProductSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	results, _ := uc.productsRepo.GetProductPaginate(ctx, payload, "buyer")

	response.WriteJSONResponse(rw, 200, results, nil)
}
