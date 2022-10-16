package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/IndominusByte/gokomodo-be/api/internal/config"
	"github.com/IndominusByte/gokomodo-be/api/internal/constant"
	authentity "github.com/IndominusByte/gokomodo-be/api/internal/entity/auth"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
)

type AuthUsecase struct {
	authRepo authRepo
}

func NewAuthUsecase(authRepo authRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, rw http.ResponseWriter, payload *authentity.JsonLoginSchema, cfg *config.Config) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	user, err := uc.authRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
			constant.App: "Invalid credential.",
		})
		return
	}

	if !uc.authRepo.IsPasswordSameAsHash(ctx, []byte(user.Password), []byte(payload.Password)) {
		response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
			constant.App: "Invalid credential.",
		})
		return
	}

	// create token
	accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: strconv.Itoa(user.Id), Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires), Fresh: true})
	refreshToken := auth.GenerateRefreshToken(&auth.RefreshToken{Sub: strconv.Itoa(user.Id), Exp: jwtauth.ExpireIn(cfg.JWT.RefreshExpires)})

	response.WriteJSONResponse(rw, 200, map[string]interface{}{
		"access_token":  auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, accessToken),
		"refresh_token": auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, refreshToken),
	}, nil)
}

func (uc *AuthUsecase) FreshToken(ctx context.Context, rw http.ResponseWriter,
	payload *authentity.JsonPasswordOnlySchema, cfg *config.Config) {

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

	if !uc.authRepo.IsPasswordSameAsHash(ctx, []byte(user.Password), []byte(payload.Password)) {
		response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
			constant.App: "Password does not match with our records.",
		})
		return
	}

	// create token
	accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: strconv.Itoa(user.Id), Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires), Fresh: true})

	response.WriteJSONResponse(rw, 200, map[string]interface{}{
		"access_token": auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, accessToken),
	}, nil)
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, rw http.ResponseWriter, cfg *config.Config) {
	_, claims, _ := jwtauth.FromContext(ctx)
	sub, _ := strconv.Atoi(claims["sub"].(string))

	user, err := uc.authRepo.GetUserById(ctx, sub)
	if err != nil {
		response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
			constant.Header: constant.UserNotFound,
		})
		return
	}

	// create token
	accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: strconv.Itoa(user.Id), Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires)})

	response.WriteJSONResponse(rw, 200, map[string]interface{}{
		"access_token": auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, accessToken),
	}, nil)
}

func (uc *AuthUsecase) AccessRevoke(ctx context.Context, rw http.ResponseWriter, redisCli *redis.Pool, cfg *config.Config) {
	conn := redisCli.Get()
	defer conn.Close()

	_, claims, _ := jwtauth.FromContext(ctx)
	conn.Do("SETEX", claims["jti"], cfg.JWT.AccessExpires.Seconds(), "ok")

	response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
		constant.App: "An access token has revoked.",
	})
}

func (uc *AuthUsecase) RefreshRevoke(ctx context.Context, rw http.ResponseWriter, redisCli *redis.Pool, cfg *config.Config) {
	conn := redisCli.Get()
	defer conn.Close()

	_, claims, _ := jwtauth.FromContext(ctx)
	conn.Do("SETEX", claims["jti"], cfg.JWT.RefreshExpires.Seconds(), "ok")

	response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
		constant.App: "An refresh token has revoked.",
	})
}
