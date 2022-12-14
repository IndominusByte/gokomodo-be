package handler_http

import (
	"net/http"
	"strings"

	"github.com/IndominusByte/gokomodo-be/api/internal/config"
	endpoint_http "github.com/IndominusByte/gokomodo-be/api/internal/endpoint/http"
	authrepo "github.com/IndominusByte/gokomodo-be/api/internal/repo/auth"
	ordersrepo "github.com/IndominusByte/gokomodo-be/api/internal/repo/orders"
	productsrepo "github.com/IndominusByte/gokomodo-be/api/internal/repo/products"
	authusecase "github.com/IndominusByte/gokomodo-be/api/internal/usecase/auth"
	ordersusecase "github.com/IndominusByte/gokomodo-be/api/internal/usecase/orders"
	productsusecase "github.com/IndominusByte/gokomodo-be/api/internal/usecase/products"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/filestatic"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Router *chi.Mux
	// Db config can be added here
	db       *sqlx.DB
	redisCli *redis.Pool
	cfg      *config.Config
}

func CreateNewServer(db *sqlx.DB, redisCli *redis.Pool, cfg *config.Config) *Server {
	s := &Server{db: db, redisCli: redisCli, cfg: cfg}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() error {
	// jwt
	publicKey, privateKey := auth.DecodeRSA(s.cfg.JWT.PublicKey, s.cfg.JWT.PrivateKey)
	TokenAuthRS256 := jwtauth.New(s.cfg.JWT.Algorithm, privateKey, publicKey)
	s.Router.Use(jwtauth.Verifier(TokenAuthRS256))

	// middleware stack
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))
	// serve file static
	fileServer := http.FileServer(filestatic.FileSystem{Static: http.Dir("static")})
	s.Router.Handle("/static/*", http.StripPrefix(strings.TrimRight("/static/", "/"), fileServer))

	// you can insert your behaviors here
	authRepo, err := authrepo.New(s.db)
	if err != nil {
		return err
	}
	authUsecase := authusecase.NewAuthUsecase(authRepo)
	endpoint_http.AddAuth(s.Router, authUsecase, s.redisCli, s.cfg)

	productsRepo, err := productsrepo.New(s.db)
	if err != nil {
		return err
	}
	productsUsecase := productsusecase.NewProductsUsecase(productsRepo, authRepo)
	endpoint_http.AddProducts(s.Router, productsUsecase, s.redisCli)

	ordersRepo, err := ordersrepo.New(s.db)
	if err != nil {
		return err
	}
	ordersUsecase := ordersusecase.NewOrdersUsecase(ordersRepo, authRepo)
	endpoint_http.AddOrders(s.Router, ordersUsecase, s.redisCli)

	return nil
}
