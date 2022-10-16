package main

import (
	"fmt"
	"log"

	_ "github.com/IndominusByte/gokomodo-be/api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/IndominusByte/gokomodo-be/api/internal/config"
	handler_http "github.com/IndominusByte/gokomodo-be/api/internal/endpoint/http/handler"
)

func startApp(cfg *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// connect the db
	db, err := config.DBConnect(cfg)
	if err != nil {
		return err
	}
	log.Printf("DB connected")

	// connect redis
	redisCli, err := config.RedisConnect(cfg)
	if err != nil {
		return err
	}
	log.Println("Redis connected")

	// mount router
	r := handler_http.CreateNewServer(db, redisCli, cfg)
	if err := r.MountHandlers(); err != nil {
		return err
	}
	log.Println("Router mounted")

	return startServer(r.Router, cfg)
}
