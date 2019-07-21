package core

import (
	"net"
	"strings"
	"erproxy/conf"
)

// ResultStatus ,
type ResultStatus int

const (
	// Direct .
	Direct ResultStatus = iota 
	// Block .
	Block	
	// Proxy .				
	Proxy
)

func getOutBound(host,port string,atype byte) outbound {

	ret := route(host,port,atype)
	if ret == Direct {
		return new(freebound)
	} else if ret == Block {
		return new(blockbound)
	} else {
		if conf.CC.OutBound.Type == "free" {
			return new(freebound)
		} else if conf.CC.OutBound.Type == "socks" {
			return new(sockbound)
		} else if conf.CC.OutBound.Type == "block" {
			return new(blockbound)
		}
	}
	return nil
}

func route(host, port string , atype byte) ResultStatus {
	for rule,policy := range(conf.CC.Routes.Route) {
		if routeMatch(host,port,atype,rule) {
			return getPolicy(policy)
		}
	}
	return getPolicy(conf.CC.Routes.Default)
}

func routeMatch(host, port string , atype byte, rule string) bool {
	testips := make([]net.IP,0)

	if atype == 0x01 || atype == 0x04 {
		t := net.ParseIP(host)
		if t != nil {
			testips = append(testips,t)
		}
	} else {
		ts,err := net.LookupHost(host)
		if err == nil {
			for _,v := range(ts) {
				t := net.ParseIP(v)
				if t != nil {
					testips = append(testips,t)
				}
			}
		}
	}

	// IP compare
	ruleIP := net.ParseIP(rule)
	if ruleIP != nil {
		for _,v := range(testips) {
			if v.Equal(ruleIP) {
				return true
			}
		}
		return false
	}

	// CIDR compare
	ruleIP,ruleCIDR,err := net.ParseCIDR(rule)
	if err == nil {
		for _,v := range(testips) {
			if ruleCIDR.Contains(v) {
				return true
			}
			return false
		}
	}

	// port compare
	if strings.Contains(rule,"port:") {
		if rule[5:] == port {
			return true
		}
	}

	// Daemon compare
	return strings.Contains(host,rule)
}

func getPolicy(v string) ResultStatus {
	if v == "direct" {
		return Direct
	} else if v == "proxy" {
		return Proxy
	} else if v == "block" {
		return Block
	}
	return Proxy
}
