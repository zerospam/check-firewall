package lib

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type TransportServer struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	OnMx   bool   `json:"mx"`
}

func (t *TransportServer) Address(server string) string {
	builder := strings.Builder{}
	builder.WriteString(server)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(t.Port))
	return builder.String()
}

func (t *TransportServer) isIp() bool {
	addr := net.ParseIP(t.Server)
	return addr != nil
}

func (t *TransportServer) getNames() (names []NameIp, err error) {

	if t.isIp() {
		names = []NameIp{
			{Name: t.Server, IP: net.ParseIP(t.Server)},
		}
		return names, nil
	}

	if t.OnMx {
		var mxRecords []*net.MX
		mxRecords, err = net.LookupMX(t.Server)
		for _, mx := range mxRecords {
			var ipRecords []net.IP
			ipRecords, err = net.LookupIP(strings.TrimRight(mx.Host, "."))
			for _, ip := range ipRecords {
				names = append(names, NameIp{Name: mx.Host, IP: ip})
			}
		}

	} else {
		var ipRecords []net.IP
		ipRecords, err = net.LookupIP(t.Server)
		for _, ip := range ipRecords {
			names = append(names, NameIp{Name: t.Server, IP: ip})
		}
	}

	if err != nil || len(names) == 0 {
		return nil, errors.New("can't find servers")
	}

	return names, nil
}

//Check if we can connect to the servers
func (t *TransportServer) CheckServer() CheckResult {
	var results []ServerResult
	var finalResult = true

	names, err := t.getNames()
	if err != nil {
		return CheckResult{Request: t, Success: false, Reason: fmt.Sprint(err)}
	}

	for _, server := range names {
		conn, err := net.DialTimeout("tcp", t.Address(server.IP.String()), 1*time.Second)
		if conn != nil {
			conn.Close()
		}
		finalResult = (err != nil) && finalResult
		results = append(results, ServerResult{NameIp: &server, Success: err != nil})
	}

	return CheckResult{Request: t, Success: finalResult, Results: results}
}
