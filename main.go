package main

import (
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/zerospam/check-firewall/lib/Handlers"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/check", Handlers.CheckTransport)
	http.HandleFunc("/healthz", Handlers.HealthCheck)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
