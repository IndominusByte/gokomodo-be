package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	ordersentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/orders"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/parser"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
)

type ordersUsecaseIface interface {
	Create(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.JsonCreateSchema)
	GetAll(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.QueryOrderSchema)
	Accept(ctx context.Context, rw http.ResponseWriter, orderId int)
}

func AddOrders(r *chi.Mux, uc ordersUsecaseIface, redisCli *redis.Pool) {
	r.Route("/orders", func(r chi.Router) {
		// protected route
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if err := auth.ValidateJWT(r.Context(), redisCli, "jwtRequired"); err != nil {
						response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
							constant.Header: err.Error(),
						})
						return
					}
					// Token is authenticated, pass it through
					next.ServeHTTP(rw, r)
				})
			})
			r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
				var p ordersentity.JsonCreateSchema

				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
						constant.Body: constant.FailedParseBody,
					})
					return
				}

				uc.Create(r.Context(), rw, &p)
			})
			r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
				var p ordersentity.QueryOrderSchema

				if err := validation.ParseRequest(&p, r.URL.Query()); err != nil {
					response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
						constant.Body: constant.FailedParseBody,
					})
					return
				}

				uc.GetAll(r.Context(), rw, &p)
			})
			r.Put("/accept/{order_id:[1-9][0-9]*}", func(rw http.ResponseWriter, r *http.Request) {
				orderId, _ := parser.ParsePathToInt("/orders/accept/(.*)", r.URL.Path)

				uc.Accept(r.Context(), rw, orderId)
			})
		})
		// public route
	})
}
