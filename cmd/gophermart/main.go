package main

import (
	"github.com/zelas91/gofermart/internal/controllers"
	"github.com/zelas91/gofermart/internal/repository"
	"github.com/zelas91/gofermart/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg := NewConfig()
	db, err := repository.NewPostgresDB(*cfg.DbURL)
	if err != nil {
		log.Fatalf("db init err : %v", err)

	}
	h := controllers.NewHandler(service.NewService(repository.NewRepository(db)))
	if err = http.ListenAndServe(*cfg.Addr, h.InitRoutes()); err != nil {
		log.Fatal(err)
	}

}
