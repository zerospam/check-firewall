package main

import (
	"CheckFirewall/lib"
	"encoding/json"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/check", checkTransport)
}

func checkTransport(w http.ResponseWriter, req *http.Request) {
	var transportServer lib.TransportServer

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	if req.Header.Get("Authorization") != os.Getenv("SHARED_KEY") {
		http.Error(w, "Wrong Key sent.", 402)
		return
	}

	json.NewDecoder(req.Body).Decode(&transportServer)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transportServer.CheckServer())

}

func main() {
	http.ListenAndServe(":8082", nil)
}
