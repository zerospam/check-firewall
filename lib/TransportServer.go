package lib

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type TransportServer struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	OnMx   bool   `json:"mx"`
}

func (t *TransportServer) Address(server string) string {
	return fmt.Sprintf("%s:%d", server, t.Port)
}

func (t *TransportServer) isIp() bool {
	addr := net.ParseIP(t.Server)
	return addr != nil
}

func (t *TransportServer) getNames() (names []NameIp, error error) {

	if t.isIp() {
		names = []NameIp{
			{Name: t.Server, IP: net.ParseIP(t.Server)},
		}
		return names, nil
	}

	if t.OnMx {
		mxRecords, errorMx := net.LookupMX(t.Server)

		if errorMx != nil {
			return nil, errorMx
		}

		names = make([]NameIp, len(mxRecords))

		for _, mx := range mxRecords {
			ipRecords, err := net.LookupIP(strings.TrimRight(mx.Host, "."))

			if err != nil {
				continue
			}

			for _, ip := range ipRecords {
				names = append(names, NameIp{Name: mx.Host, IP: ip})
			}
		}

	} else {
		ipRecords, errIp := net.LookupIP(t.Server)

		if errIp != nil {
			return nil, errIp
		}

		names = make([]NameIp, 1)

		for _, ip := range ipRecords {
			names = append(names, NameIp{Name: t.Server, IP: ip})
		}
	}

	if len(names) == 0 {
		return nil, errors.New("can't find servers")
	}

	return names, nil
}

//Check if we can connect to the servers
func (t *TransportServer) CheckServer() CheckResult {
	names, err := t.getNames()

	if err != nil {
		return CheckResult{Request: t, Success: false, Reason: fmt.Sprint(err)}
	}

	var finalResult = true
	var results = make([]ServerResult, len(names))

	for index, server := range names {
		conn, err := net.DialTimeout("tcp", t.Address(server.IP.String()), 1*time.Second)
		if conn != nil {
			conn.Close()
		}
		finalResult = (err != nil) && finalResult
		results[index] = ServerResult{NameIp: &server, Success: err != nil}
	}

	return CheckResult{Request: t, Success: finalResult, Results: results}
}
