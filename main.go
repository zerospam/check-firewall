package main

import (
	"fmt"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/zerospam/check-firewall/lib/Handlers"
	"github.com/zerospam/check-firewall/lib/common"
	"net/http"
)

func init() {
	http.HandleFunc("/check", Handlers.CheckTransport)
	http.HandleFunc("/healthz", Handlers.HealthCheck)
}

func main() {

	err := http.ListenAndServe(fmt.Sprintf(":%s", common.GetVars().ApplicationPort), nil)
	if err != nil {
		panic(err)
	}
}
