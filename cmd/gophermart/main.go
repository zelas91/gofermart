package main

import (
	"github.com/zelas91/gofermart/internal/router"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", router.InitRoutes())
}
