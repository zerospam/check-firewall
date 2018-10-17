package main

import (
	"encoding/json"
	"fmt"
	"github.com/zerospam/check-firewall/lib"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	http.HandleFunc("/check", checkTransport)
}

func checkTransport(w http.ResponseWriter, req *http.Request) {
	var transportServer lib.TransportServer

	if req.Header.Get("Authorization") != os.Getenv("SHARED_KEY") {
		http.Error(w, "Wrong Key sent.", 402)
		log.Printf("[%s] - %s (%s:%d) - %v\n", req.RemoteAddr, req.Method, transportServer.Server, transportServer.Port, "REJECT")
		return
	}

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	json.NewDecoder(req.Body).Decode(&transportServer)

	w.Header().Add("Content-Type", "application/json")
	result := transportServer.CheckServer()
	json.NewEncoder(w).Encode(result)
	log.Printf("[%s] - %s (%s:%d) - %v\n", req.RemoteAddr, req.Method, transportServer.Server, transportServer.Port, result.Success)

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
