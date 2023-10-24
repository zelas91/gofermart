package main

import (
	"context"
	"errors"
	"github.com/zelas91/gofermart/internal/controllers"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/repository"
	"github.com/zelas91/gofermart/internal/service"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	http *http.Server
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel

	cfg := NewConfig()
	log := logger.New()
	db, err := repository.NewPostgresDB(*cfg.DBURL)
	if err != nil {
		log.Fatalf("db init err : %v", err)

	}

	log.Info("Accrual add ", *cfg.Accrual)
	h := controllers.NewHandler(service.NewService(repository.NewRepository(db)))
	serv := &server{http: &http.Server{Addr: *cfg.Addr, Handler: h.InitRoutes(log)}}
	go func() {
		if err = serv.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()

	log.Infof("start server, address : %s", *cfg.Addr)

	<-ctx.Done()

	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err = serv.http.Shutdown(ctxTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("shutdown server %v", err)
	}

	if err = db.Close(); err != nil {
		log.Error(err)
	}

	log.Info("server stop")
}
