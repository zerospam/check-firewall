package lib

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/zerospam/check-firewall/lib/common"
	"github.com/zerospam/check-firewall/lib/tlsgenerator"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type TransportServer struct {
	Server       string `json:"server"`
	Port         int    `json:"port"`
	OnMx         bool   `json:"mx"`
	TestEmail    string `json:"test_email"`
	tlsGenerator *tlsgenerator.CertificateGenerator
}

func (t *TransportServer) Address(server string) string {
	return fmt.Sprintf("%s:%d", server, t.Port)
}

func (t *TransportServer) isIp() bool {
	addr := net.ParseIP(t.Server)
	return addr != nil
}

func (t *TransportServer) getClientTLSConfig(commonName string) *tls.Config {
	if t.tlsGenerator == nil {
		t.tlsGenerator = tlsgenerator.NewClient(time.Now(), 365*24*time.Hour)
	}

	return t.tlsGenerator.GetTlsClientConfig(commonName)
}

func (t *TransportServer) getNames() (names []*NameIp, error error) {

	if t.isIp() {
		names = []*NameIp{
			{Name: t.Server, IP: net.ParseIP(t.Server)},
		}
		return names, nil
	}

	if t.OnMx {
		mxRecords, errorMx := net.LookupMX(t.Server)

		if errorMx != nil {
			return nil, errorMx
		}

		for _, mx := range mxRecords {
			ipRecords, err := net.LookupIP(strings.TrimRight(mx.Host, "."))

			if err != nil {
				continue
			}

			for _, ip := range ipRecords {
				names = append(names, &NameIp{Name: mx.Host, IP: ip})
			}
		}

	} else {
		ipRecords, errIp := net.LookupIP(t.Server)

		if errIp != nil {
			return nil, errIp
		}

		for _, ip := range ipRecords {
			names = append(names, &NameIp{Name: t.Server, IP: ip})
		}
	}

	if len(names) == 0 {
		return nil, errors.New("can't find servers")
	}

	return names, nil
}

//Check if we can connect to the servers
func (t *TransportServer) CheckServer(checkSMTP bool) CheckResult {
	names, err := t.getNames()

	if err != nil {
		return CheckResult{Request: t, Success: false, Reason: fmt.Sprint(err)}
	}

	var finalResult = true
	var results = make([]ServerResult, len(names))

	for index, server := range names {
		conn, err := net.DialTimeout("tcp", t.Address(server.IP.String()), 1*time.Second)
		currentSuccess := err != nil
		var msg string
		if conn != nil {
			if checkSMTP {
				currentSuccess, msg = t.checkSMTP(conn)
			}
			conn.Close()
		}
		finalResult = currentSuccess && finalResult
		results[index] = ServerResult{NameIp: server, Success: currentSuccess, Message: msg}
	}

	return CheckResult{Request: t, Success: finalResult, Results: results}
}

func (t *TransportServer) checkSMTP(conn net.Conn) (bool, string) {

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	client, err := smtp.NewClient(conn, t.Server)
	if err != nil {
		return true, "Stop at INIT CONN"
	}

	defer client.Quit()
	defer client.Close()

	if err = client.Hello(hostname); err != nil {
		return true, "Stop at EHLO"
	}

	if tlsSupport, _ := client.Extension("STARTTLS"); tlsSupport {
		tlsConfig := t.getClientTLSConfig(common.GetVars().SmtpCN)
		tlsConfig.ServerName = t.Server
		err = client.StartTLS(tlsConfig)
		if err != nil {
			log.Printf("Couldn't start TLS transaction: %s", err)
			return true, fmt.Sprintf("Couldn't start TLS transaction: %s", err)
		}
	}

	if err = client.Mail(common.GetVars().SmtpMailFrom.String()); err != nil {
		return true, "Stop at MAIL FROM"
	}

	if err = client.Rcpt(t.TestEmail); err != nil {
		return true, "Stop at RCPT TO"
	}

	return false, "Can start a SMTP Transaction"
}
