package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/gokomodo-be/api/internal/config"
	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/response"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
)

type authUsecaseIface interface {
	Login(ctx context.Context, rw http.ResponseWriter, payload *authentity.JsonLoginSchema, cfg *config.Config)
	FreshToken(ctx context.Context, rw http.ResponseWriter, payload *authentity.JsonPasswordOnlySchema, cfg *config.Config)
	RefreshToken(ctx context.Context, rw http.ResponseWriter, cfg *config.Config)
	AccessRevoke(ctx context.Context, rw http.ResponseWriter, redisCli *redis.Pool, cfg *config.Config)
	RefreshRevoke(ctx context.Context, rw http.ResponseWriter, redisCli *redis.Pool, cfg *config.Config)
}

func AddAuth(r *chi.Mux, uc authUsecaseIface, redisCli *redis.Pool, cfg *config.Config) {
	r.Route("/auth", func(r chi.Router) {
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
			r.Post("/fresh-token", func(rw http.ResponseWriter, r *http.Request) {
				var p authentity.JsonPasswordOnlySchema

				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
						constant.Body: constant.FailedParseBody,
					})
					return
				}

				uc.FreshToken(r.Context(), rw, &p, cfg)
			})
			r.Delete("/access-revoke", func(rw http.ResponseWriter, r *http.Request) {
				uc.AccessRevoke(r.Context(), rw, redisCli, cfg)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if err := auth.ValidateJWT(r.Context(), redisCli, "jwtRefreshRequired"); err != nil {
						response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
							constant.Header: err.Error(),
						})
						return
					}
					// Token is authenticated, pass it through
					next.ServeHTTP(rw, r)
				})
			})
			r.Post("/refresh-token", func(rw http.ResponseWriter, r *http.Request) {
				uc.RefreshToken(r.Context(), rw, cfg)
			})
			r.Delete("/refresh-revoke", func(rw http.ResponseWriter, r *http.Request) {
				uc.RefreshRevoke(r.Context(), rw, redisCli, cfg)
			})
		})

		// public route
		r.Post("/login", func(rw http.ResponseWriter, r *http.Request) {
			var p authentity.JsonLoginSchema

			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					constant.Body: constant.FailedParseBody,
				})
				return
			}

			uc.Login(r.Context(), rw, &p, cfg)
		})
	})
}
