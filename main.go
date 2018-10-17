package main

import (
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/zerospam/check-firewall/lib/Handlers"
	"net/http"
	"os"
	"strconv"
)

func init() {
	http.HandleFunc("/check", Handlers.CheckTransport)
	http.HandleFunc("/healthz", Handlers.HealthCheck)
}

func main() {
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
	if port == 0 || err != nil {
		port = 80
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
