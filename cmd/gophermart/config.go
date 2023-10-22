package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var (
	addr    *string
	dbURL   *string
	accrual *string
)

func init() {
	addr = flag.String("a", "localhost:8081", "endpoint start server")
	dbURL = flag.String("d", "host=localhost port=5432 user=userm dbname=gofermart password=12345678 sslmode=disable", "url DB")
	accrual = flag.String("r", "", "Database URL")
}

type Config struct {
	Addr    *string `env:"RUN_ADDRESS"`
	DBURL   *string `env:"DATABASE_URI"`
	Accrual *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("read env error=%v", err)
	}

	if cfg.Addr == nil {
		cfg.Addr = addr
	}
	if cfg.DBURL == nil {
		cfg.DBURL = dbURL
	}
	if cfg.Accrual == nil {
		cfg.Accrual = accrual
	}
	flag.Parse()
	return &cfg
}
