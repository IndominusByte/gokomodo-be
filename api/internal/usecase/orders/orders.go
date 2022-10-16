package orders

import (
	"context"
	"net/http"
	"strconv"

	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	ordersentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/orders"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
	"github.com/go-chi/jwtauth"
)

type OrdersUsecase struct {
	ordersRepo ordersRepo
	authRepo   authRepo
}

func NewOrdersUsecase(orderRepo ordersRepo, authRepo authRepo) *OrdersUsecase {
	return &OrdersUsecase{
		ordersRepo: orderRepo,
		authRepo:   authRepo,
	}
}

func (uc *OrdersUsecase) Create(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.JsonCreateSchema) {
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

	product, err := uc.ordersRepo.GetProductById(ctx, payload.ProductId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Product not found.",
		})
		return
	}

	if user.Id == product.SellerId {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			constant.App: "Cannot buy your own product.",
		})
		return
	}

	// insert into db
	payload.BuyerId = user.Id
	payload.SellerId = product.SellerId
	payload.ProductId = product.Id
	payload.SourceAddress = product.SellerAddress
	payload.DestinationAddress = user.Address
	payload.Price = product.Price
	payload.TotalPrice = payload.Qty * product.Price
	uc.ordersRepo.Insert(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		constant.App: "Successfully add a new order.",
	})

}

func (uc *OrdersUsecase) GetAll(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.QueryOrderSchema) {
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

	switch payload.For {
	case "buyer":
		payload.BuyerId = user.Id
	case "seller":
		payload.SellerId = user.Id
	}

	results, _ := uc.ordersRepo.GetOrderPaginate(ctx, payload)

	response.WriteJSONResponse(rw, 200, results, nil)
}

func (uc *OrdersUsecase) Accept(ctx context.Context, rw http.ResponseWriter, orderId int) {
	_, claims, _ := jwtauth.FromContext(ctx)
	sub, _ := strconv.Atoi(claims["sub"].(string))

	user, err := uc.authRepo.GetUserById(ctx, sub)
	if err != nil {
		response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
			constant.Header: constant.UserNotFound,
		})
		return
	}

	_, err = uc.ordersRepo.GetOrderById(ctx, orderId, user.Id)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Order not found.",
		})
		return

	}

	// set accepted
	uc.ordersRepo.SetAcceptedOrder(ctx, orderId)

	response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
		constant.App: "Successfully accept the order.",
	})
}
