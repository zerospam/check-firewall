package main

import (
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/zerospam/check-firewall/lib/Common"
	"github.com/zerospam/check-firewall/lib/Handlers"
	"net/http"
)

func init() {
	http.HandleFunc("/check", Handlers.CheckTransport)
	http.HandleFunc("/healthz", Handlers.HealthCheck)
}

func main() {

	err := http.ListenAndServe(fmt.Sprintf(":%s", Common.GetVars().ApplicationPort), nil)
	if err != nil {
		panic(err)
	}
}
