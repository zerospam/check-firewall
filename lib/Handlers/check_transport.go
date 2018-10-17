package Handlers

import (
	"encoding/json"
	"github.com/zerospam/check-firewall/lib"
	"log"
	"net/http"
	"os"
	"strings"
)

func getRequestIp(req *http.Request) string {
	if header := req.Header.Get("X-Forwarded-For"); header != "" {
		exploded := strings.Split(header, ",")
		return strings.Trim(exploded[len(exploded)-1], " ")
	}

	return req.RemoteAddr
}

func CheckTransport(w http.ResponseWriter, req *http.Request) {
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

	defer req.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	result := transportServer.CheckServer()
	json.NewEncoder(w).Encode(result)
	log.Printf("[%s] - %s (%s:%d) - %v\n", getRequestIp(req), req.Method, transportServer.Server, transportServer.Port, result.Success)

}
