package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	productsentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/products"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
)

type productsUsecaseIface interface {
	Create(ctx context.Context, rw http.ResponseWriter, payload *productsentity.JsonCreateSchema)
	MyProduct(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryProductSchema)
	GetAll(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryProductSchema)
}

func AddProducts(r *chi.Mux, uc productsUsecaseIface, redisCli *redis.Pool) {
	r.Route("/products", func(r chi.Router) {
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
				var p productsentity.JsonCreateSchema

				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
						constant.Body: constant.FailedParseBody,
					})
					return
				}

				uc.Create(r.Context(), rw, &p)
			})
			r.Get("/mine", func(rw http.ResponseWriter, r *http.Request) {
				var p productsentity.QueryProductSchema

				if err := validation.ParseRequest(&p, r.URL.Query()); err != nil {
					response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
						constant.Body: constant.FailedParseBody,
					})
					return
				}

				uc.MyProduct(r.Context(), rw, &p)
			})
		})
		// public route
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			var p productsentity.QueryProductSchema

			if err := validation.ParseRequest(&p, r.URL.Query()); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					constant.Body: constant.FailedParseBody,
				})
				return
			}

			uc.GetAll(r.Context(), rw, &p)
		})
	})
}
