package main

import (
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/zerospam/check-firewall/lib/environment-vars"
	"github.com/zerospam/check-firewall/lib/handlers"
	"net/http"
)

func init() {
	http.HandleFunc("/check", handlers.CheckTransport)
	http.HandleFunc("/healthz", handlers.HealthCheck)
}

func main() {

	err := http.ListenAndServe(fmt.Sprintf(":%s", environmentvars.GetVars().ApplicationPort), nil)
	if err != nil {
		panic(err)
	}
}
