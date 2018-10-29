package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerospam/check-firewall/lib"
	"github.com/zerospam/check-firewall/lib/Handlers"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ListenFreePort() (*net.TCPListener, int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return nil, 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, 0, err
	}
	return l, l.Addr().(*net.TCPAddr).Port, nil
}

func TestHealthz(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handlers.HealthCheck)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFirewalled(t *testing.T) {
	port := 1
	server := lib.TransportServer{Server: "localhost", Port: port, OnMx: false}

	result := server.CheckServer(false)

	assert.True(t, result.Success, "Expected to have localhost:%d unreachable", port)
	assert.True(t, result.Results[0].Success, "Expected to have localhost:%d unreachable", port)
	assert.Equal(t, net.ParseIP("127.0.0.1"), result.Results[0].IP, "Expected to have localhost:%d unreachable", port)
}

func TestNotFirewalled(t *testing.T) {
	listener, port, err := ListenFreePort()
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	server := lib.TransportServer{Server: "localhost", Port: port, OnMx: false}
	result := server.CheckServer(false)

	assert.False(t, result.Success, "Expected to have localhost:%d reachable", port)
	assert.False(t, result.Results[0].Success, "Expected to have localhost:%d reachable", port)
	assert.Equal(t, net.ParseIP("127.0.0.1"), result.Results[0].IP, "Expected to have localhost:%d unreachable", port)

}
