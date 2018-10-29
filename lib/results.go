package lib

import "net"

type NameIp struct {
	Name string `json:"server"`
	IP   net.IP `json:"ip"`
}

type ServerResult struct {
	*NameIp
	Success bool   `json:"result"`
	Message string `json:"smtp_check_msg"`
}
type CheckResult struct {
	Request *TransportServer `json:"request"`
	Success bool             `json:"success"`
	Reason  string           `json:"error_message,omitempty"`
	Results []ServerResult   `json:"details,omitempty"`
}
